package domain

import "time"

type Order struct {
	ID        int         `json:"id" db:"id"`
	UserID    int         `json:"user_id" db:"user_id"`
	Status    string      `json:"status" db:"status"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	Items     []OrderItem `json:"items"`
}

type OrderItem struct {
	ID        int `json:"id" db:"id"`
	OrderID   int `json:"order_id" db:"order_id"`
	ProductID int `json:"product_id" db:"product_id"`
	Quantity  int `json:"quantity" db:"quantity"`
}

type OrderRepository interface {
	Create(order *Order) error
	GetByID(id int) (*Order, error)
	UpdateStatus(id int, status string) error
	ListByUser(userID int) ([]Order, error)
}

type OrderUsecase interface {
	Create(order *Order) error
	GetByID(id int) (*Order, error)
	UpdateStatus(id int, status string) error
	ListByUser(userID int) ([]Order, error)
}
