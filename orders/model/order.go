package model

import "time"

type Order struct {
	Id    uint64      `json:"id"`    // order id
	Date  time.Time   `json:"date"`  // date of order
	Total float64     `json:"total"` // total price
	Items []OrderItem `json:"items"` // order items
	// TODO: add user module
	// CustomerId uint64  `json:"customer_id"`
}

type OrderItem struct {
	Id       uint64  `json:"id"`                 // PK
	OrderId  uint64  `json:"order_id,omitempty"` // FK to order
	BookId   uint64  `json:"book_id,omitempty"`  // FK to book
	Quantity uint8   `json:"quantity"`           // quantity
	Price    float64 `json:"price"`              // price
	Author   string  `json:"author"`
	Cover    string  `json:"cover"`
	Name     string  `json:"name"`
	// TODO: add user module
	// CustomerId uint64  `json:"customer_id"`
}

type CheckoutOrder struct {
	Date  time.Time `json:"date"`  // date of order
	Total float64   `json:"total"` // total price
}

type CheckoutItem struct {
	BookId   uint64  `json:"book_id"`
	Quantity uint8   `json:"quantity"`
	Price    float64 `json:"price"`
}

type OrderPlaced struct {
	Books []uint64 `json:"books"`
}
