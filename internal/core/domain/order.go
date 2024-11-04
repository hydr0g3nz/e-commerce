package domain

import (
	"errors"
	"time"
)

type Order struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	Date            time.Time `json:"date"`
	Status          string    `json:"status"`
	ShippingAddress Address   `json:"shipping_address"`
	Items           []Item    `json:"items"`
	TotalPrice      float64   `json:"total_price"`
	PaymentMethod   string    `json:"payment_method"`
}
type Item struct {
	Id       string  `json:"product_id"`
	Sku      string  `json:"sku"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
	Sale     int     `json:"sale"`
}

func (o *Order) ValidateCreate() error {
	if o.UserID == "" {
		return errors.New("user id is required")
	}
	if o.Status == "" {
		return errors.New("status is required")
	}
	if o.ShippingAddress == (Address{}) {
		return errors.New("shipping address is required")
	}
	if len(o.Items) == 0 {
		return errors.New("items are required")
	}
	if o.PaymentMethod == "" {
		return errors.New("payment method is required")
	}
	return nil
}
