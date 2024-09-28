package domain

import "time"

type Review struct {
	// "review_id": "rev001",
	// "product_id": "987",
	// "user_id": "123",
	// "rating": 4,
	// "comment": "Great mouse, works smoothly.",
	// "date": "2024-09-02"
	ID        string    `json:"review_id"`
	ProductID string    `json:"product_id"`
	UserID    string    `json:"user_id"`
	Rating    float64   `json:"rating"`
	Comment   string    `json:"comment"`
	Date      time.Time `json:"date"`
}
