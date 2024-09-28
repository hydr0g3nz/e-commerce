package domain

import "time"

// {
// 	"order_id": "ord001",
// 	"user_id": "123",
// 	"date": "2024-09-01",
// 	"status": "Shipped",
// 	"shipping_address": {
// 	  "street": "123 Main St",
// 	  "city": "Metropolis",
// 	  "state": "NY",
// 	  "zip": "10001"
// 	},
// 	"items": [
// 	  {
// 		"product_id": "987",
// 		"name": "Wireless Mouse",
// 		"quantity": 2,
// 		"price": 25.00
// 	  },
// 	  {
// 		"product_id": "654",
// 		"name": "Laptop",
// 		"quantity": 1,
// 		"price": 1200.00
// 	  }
// 	],
// 	"total_price": 1250.00,
// 	"payment_method": "Credit Card"
//   }

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
