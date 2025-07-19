package models

type Product struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	Image      string  `json:"image"`
	Price      float64 `json:"price"`
	Qty        int     `json:"qty"`
	OutOfStock bool    `json:"out_of_stock"`
}
