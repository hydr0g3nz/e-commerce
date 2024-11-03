package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hydr0g3nz/e-commerce/internal/core/domain"
	"github.com/hydr0g3nz/e-commerce/internal/core/ports"
	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	OrderStatusPending    = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusCompleted  = "completed"
	OrderStatusFailed     = "failed"
	OrderStatusCancelled  = "cancelled"
	ReservationTimeout    = 15 * time.Minute
	ReservationQueueName  = "product_reservations"
	ReservationRoutingKey = "product.reserve"
	ReservationExchange   = "order_events"
)

type ReservationMessage struct {
	OrderID   string        `json:"order_id"`
	Items     []domain.Item `json:"items"`
	Timestamp time.Time     `json:"timestamp"`
}

type OrderService struct {
	orderRepo      ports.OrderRepository
	productRepo    ports.ProductRepository
	amqpChannel    *amqp.Channel
	amqpConnection *amqp.Connection
}

func NewOrderService(
	orderRepo ports.OrderRepository,
	productRepo ports.ProductRepository,
	amqpURL string,
) (*OrderService, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare exchange
	err = ch.ExchangeDeclare(
		ReservationExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Declare queue
	_, err = ch.QueueDeclare(
		ReservationQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		ReservationQueueName,
		ReservationRoutingKey,
		ReservationExchange,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &OrderService{
		orderRepo:      orderRepo,
		productRepo:    productRepo,
		amqpChannel:    ch,
		amqpConnection: conn,
	}, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) error {
	// Generate order
	order.Status = OrderStatusPending
	// Validate items and calculate total price
	if err := s.validateAndCalculateOrder(ctx, order); err != nil {
		return err
	}
	// Try to reserve products
	if err := s.sendReservationRequest(order); err != nil {
		return err
	}
	// Save order
	if err := s.orderRepo.Create(ctx, order); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) validateAndCalculateOrder(ctx context.Context, order *domain.Order) error {
	var totalPrice float64

	for i, item := range order.Items {
		// Fetch product variation to validate availability and price
		fmt.Println("id", item.Id, "sku", item.Sku)
		product, err := s.productRepo.GetProductBySku(ctx, item.Id, item.Sku)
		if err != nil {
			return err
		}

		var variation *domain.Variation
		for _, v := range product.Variations {
			if v.Sku == item.Sku {
				variation = &v
				break
			}
		}

		if variation == nil {
			return errors.New("product variation not found")
		}

		// Validate stock
		if variation.Stock < item.Quantity {
			return errors.New("insufficient stock")
		}

		// Calculate price with sale if applicable
		finalPrice := variation.Price
		if variation.Sale > 0 {
			finalPrice = variation.Price * (1 - float64(variation.Sale)/100)
		}

		// Update item with current price and sale
		order.Items[i].Price = int(finalPrice)
		order.Items[i].Sale = int(variation.Sale)

		totalPrice += float64(item.Quantity) * finalPrice
	}

	order.TotalPrice = totalPrice
	return nil
}

func (s *OrderService) sendReservationRequest(order *domain.Order) error {
	msg := ReservationMessage{
		OrderID:   order.ID,
		Items:     order.Items,
		Timestamp: time.Now(),
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	fmt.Println("send reservation request", string(body))
	return s.amqpChannel.PublishWithContext(
		context.Background(),
		ReservationExchange,
		ReservationRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Expiration:  (ReservationTimeout).String(),
		},
	)
}

func (s *OrderService) StartReservationConsumer() error {
	msgs, err := s.amqpChannel.Consume(
		ReservationQueueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var reservation ReservationMessage
			if err := json.Unmarshal(msg.Body, &reservation); err != nil {
				msg.Nack(false, true)
				continue
			}

			// Process reservation
			if err := s.processReservation(context.Background(), &reservation); err != nil {
				msg.Nack(false, true)
				continue
			}

			msg.Ack(false)
		}
	}()

	return nil
}

func (s *OrderService) processReservation(ctx context.Context, msg *ReservationMessage) error {
	fmt.Println("process reservation", msg)
	// Update order status to processing
	if err := s.orderRepo.UpdateStatus(ctx, msg.OrderID, OrderStatusProcessing); err != nil {
		return err
	}

	// Reserve inventory for each item
	for _, item := range msg.Items {
		if err := s.productRepo.ReserveStock(ctx, item.Id, item.Sku, item.Quantity); err != nil {
			// If reservation fails, release all previous reservations
			s.rollbackReservations(ctx, msg.OrderID, msg.Items)
			s.orderRepo.UpdateStatus(ctx, msg.OrderID, OrderStatusFailed)
			return err
		}
	}

	return nil
}

func (s *OrderService) rollbackReservations(ctx context.Context, orderID string, items []domain.Item) {
	for _, item := range items {
		s.productRepo.ReleaseStock(ctx, item.Id, item.Sku, item.Quantity)
	}
}

func (s *OrderService) Close() {
	if s.amqpChannel != nil {
		s.amqpChannel.Close()
	}
	if s.amqpConnection != nil {
		s.amqpConnection.Close()
	}
}
