package model

type CartItem struct {
	Id       uint64  `json:"id"`      // PK
	BookId   uint64  `json:"book_id"` // FK to book
	Quantity uint8   `json:"quantity"`
	Price    float64 `json:"price"`
	Author   string  `json:"author"`
	Cover    string  `json:"cover"`
	Name     string  `json:"name"`
	// TODO: add user module
	// CustomerId uint64  `json:"customer_id"`
}

type AddCart struct {
	BookId   uint64  `json:"book_id"`
	Quantity uint8   `json:"quantity"`
	Price    float64 `json:"price"`
}

type UpdateCart struct {
	Quantity uint8 `json:"quantity"`
}

type OrderPlaced struct {
	Books []uint64 `json:"books"`
}
