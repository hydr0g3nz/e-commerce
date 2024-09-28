package domain

// "product_id": "987",
// "name": "Wireless Mouse",
// "quantity": 2,
// "price": 25.00
type Item struct {
	Id       string `json:"product_id"`
	Sku      string `json:"Sku"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type Cart struct {
	Items      []Item  `json:"items"`
	TotalPrice float64 `json:"total_price"`
}
type Wishlist struct {
	Items []string `json:"items"`
}
