package models

import (
	"time"
)

type OrderRequest struct {
	Pair     string  `json:"pair" binding:"required"`
	Side     string  `json:"side" binding:"required"`
	Price    float64 `json:"price" binding:"required,gt=0"`
	Quantity float64 `json:"quantity" binding:"required,gt=0"`
	UserID   int     `json:"user_id" binding:"required"`
}

type Order struct {
	ID             int       `json:"id" db:"id"`
	UserID         int       `json:"user_id" db:"user_id"`
	Pair           string    `json:"pair" db:"pair"`
	Side           string    `json:"side" db:"side"`
	Price          float64   `json:"price" db:"price"`
	Quantity       float64   `json:"quantity" db:"quantity"`
	FilledQuantity float64   `json:"filled_quantity" db:"filled_quantity"`
	Status         string    `json:"status" db:"status"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type OrderResponse struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	Pair           string    `json:"pair"`
	Side           string    `json:"side"`
	Price          float64   `json:"price"`
	Quantity       float64   `json:"quantity"`
	FilledQuantity float64   `json:"filled_quantity"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// ToOrder converts OrderRequest to Order
func (or *OrderRequest) ToOrder() *Order {
	return &Order{
		UserID:         or.UserID,
		Pair:           or.Pair,
		Side:           or.Side,
		Price:          or.Price,
		Quantity:       or.Quantity,
		FilledQuantity: 0,
		Status:         "open",
		CreatedAt:      time.Now(),
	}
}
