package model

type Book struct {
	Id       uint64  `json:"id"`
	Author   string  `json:"author"`
	Category string  `json:"category"`
	Cover    string  `json:"cover"` // url to file
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Rating   float64 `json:"rating"`
}

type BestSeller struct {
	Id    uint64  `json:"id"`
	Sales float64 `json:"sales"`
}
