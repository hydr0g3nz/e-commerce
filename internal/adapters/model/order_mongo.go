package model

import "github.com/hydr0g3nz/e-commerce/internal/core/domain"

type Order struct {
	Model           `bson:"inline"`
	UserID          string         `json:"user_id" bson:"user_id"`
	Status          string         `json:"status" bson:"status"`
	ShippingAddress domain.Address `json:"shipping_address" bson:"shipping_address"`
	Items           []domain.Item  `json:"items" bson:"items"`
	TotalPrice      float64        `json:"total_price" bson:"total_price"`
	PaymentMethod   string         `json:"payment_method" bson:"payment_method"`
}

func DomainOrderToModel(o *domain.Order) *Order {
	return &Order{
		Model:           Model{ID: o.ID},
		UserID:          o.UserID,
		Status:          o.Status,
		ShippingAddress: o.ShippingAddress,
		Items:           o.Items,
		TotalPrice:      o.TotalPrice,
		PaymentMethod:   o.PaymentMethod,
	}
}
func (o *Order) ToDomain() *domain.Order {
	return &domain.Order{
		ID:              o.ID,
		UserID:          o.UserID,
		Status:          o.Status,
		Date:            o.CreatedAt,
		ShippingAddress: o.ShippingAddress,
		Items:           o.Items,
		TotalPrice:      o.TotalPrice,
		PaymentMethod:   o.PaymentMethod,
	}
}
func OrdersModelToDomainList(orders []*Order) []*domain.Order {
	var ordersList []*domain.Order
	for _, order := range orders {
		ordersList = append(ordersList, order.ToDomain())
	}
	return ordersList
}
