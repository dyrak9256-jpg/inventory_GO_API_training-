package product

import (
	"errors"
	"time"
)

var (
	ErrNotFound    = errors.New("product not found")
	ErrOutOfStock  = errors.New("not enough items in stock")
	ErrInvalidData = errors.New("invalid product data")
)

type Product struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatedProductRequest struct {
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
}

type UpdateProductRequest struct {
	Name     *string  `json:"name"`
	Category *string  `json:"category"`
	Price    *float64 `json:"price"`
	Stock    *int     `json:"stock"`
}

type BuyProductRequest struct {
	Quantity int `json:"quantity"`
}
