package domain

// "user_id": "123",
//   "name": "John Doe",
//   "email": "john@example.com",

type User struct {
	ID      string  `json:"user_id"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Address Address `json:"address"`
}
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip"`
}
