package domain

import "time"

type Review struct {
	ID        string    `json:"review_id"`
	ProductID string    `json:"product_id"`
	UserID    string    `json:"user_id"`
	Rating    float64   `json:"rating"`
	Comment   string    `json:"comment"`
	Date      time.Time `json:"date"`
}
