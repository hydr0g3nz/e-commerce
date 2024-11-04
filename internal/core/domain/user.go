package domain

// "user_id": "123",
//   "name": "John Doe",
//   "email": "john@example.com",

type User struct {
	ID       string    `json:"user_id"`
	Password string    `json:"password,omitempty" validate:"required,min=8"`
	Name     string    `json:"name" validate:"required,min=2"`
	Email    string    `json:"email" validate:"required,email"`
	Role     string    `json:"role"`
	Address  []Address `json:"address"`
}
type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip"`
}
