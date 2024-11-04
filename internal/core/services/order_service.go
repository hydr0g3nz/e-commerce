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
	ReservationQueueName  = "product.reserve"
	ReservationRoutingKey = "product.reserve"
	ReservationExchange   = "order_events"
)

type ReservationMessage struct {
	OrderID    string        `json:"order_id"`
	Items      []domain.Item `json:"items"`
	Timestamp  time.Time     `json:"timestamp"`
	RetryCount int           `json:"retry_count"`
}

type OrderService struct {
	orderRepo      ports.OrderRepository
	productRepo    ports.ProductRepository
	amqpChannel    *amqp.Channel
	amqpConnection *amqp.Connection
	amqpQueueName  string
}

func NewOrderService(
	orderRepo ports.OrderRepository,
	productRepo ports.ProductRepository,
	amqpURL string,
) (*OrderService, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		fmt.Println("Error connecting to RabbitMQ:", err)
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("Error creating channel:", err)
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
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		false,
		false,
		amqp.Table{
			"x-dead-letter-exchange": "ecom_dlx",
		},
	)
	if err != nil {
		return nil, err
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		q.Name,
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
		amqpQueueName:  q.Name,
	}, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) error {
	// Generate order
	order.Status = OrderStatusPending
	// Validate items and calculate total price
	if err := s.validateAndCalculateOrder(ctx, order); err != nil {
		fmt.Println("Error validating and calculating order:", err)
		return err
	}
	// Save order
	orderId, err := s.orderRepo.Create(ctx, order)
	if err != nil {
		return err
	}
	order.ID = orderId
	// Try to reserve products
	if err := s.sendReservationRequest(order); err != nil {
		fmt.Println("Error sending reservation request:", err)
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
		OrderID:    order.ID,
		Items:      order.Items,
		Timestamp:  time.Now(),
		RetryCount: 0,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return s.amqpChannel.PublishWithContext(
		context.Background(),
		ReservationExchange,
		ReservationRoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			// Expiration:  (ReservationTimeout).String(),
		},
	)
}

func (s *OrderService) StartReservationConsumer() error {
	msgs, err := s.amqpChannel.Consume(
		s.amqpQueueName,
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
				fmt.Println("error unmarshalling reservation", err)
				msg.Nack(false, true)
				continue
			}
			if reservation.RetryCount >= 3 {
				msg.Reject(false)
				continue
			}
			// Process reservation
			if err := s.processReservation(context.Background(), &reservation); err != nil {
				fmt.Println("error processing reservation", err)
				reservation.RetryCount++
				msgBody, err := json.Marshal(reservation)
				if err != nil {
					fmt.Println("error marshalling reservation", err)
					continue
				}
				err = s.amqpChannel.PublishWithContext(
					context.Background(),
					ReservationExchange,
					ReservationRoutingKey,
					false,
					false,
					amqp.Publishing{
						ContentType: "application/json",
						Body:        msgBody,
					},
				)
				if err != nil {
					fmt.Println("error publishing reservation", err)
					continue
				}
				msg.Ack(false)
				continue
			}
			msg.Ack(false)
		}
	}()

	return nil
}

func (s *OrderService) processReservation(ctx context.Context, msg *ReservationMessage) error {
	// Update order status to processing
	if err := s.orderRepo.UpdateStatus(ctx, msg.OrderID, OrderStatusProcessing); err != nil {
		return err
	}

	// Reserve inventory for each item
	for _, item := range msg.Items {
		if err := s.productRepo.ReserveStock(ctx, item.Id, item.Sku, item.Quantity); err != nil {
			// If reservation fails, release all previous reservations
			fmt.Println("reservation failed", err)
			s.rollbackReservations(ctx, msg.OrderID, msg.Items)
			if err := s.orderRepo.UpdateStatus(ctx, msg.OrderID, OrderStatusFailed); err != nil {
				return err
			}
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
