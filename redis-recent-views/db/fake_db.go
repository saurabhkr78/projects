package db

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var Products = []Product{
	{ID: "1", Name: "RAM", Price: 1002.0},
	{ID: "2", Name: "CDs", Price: 1003.0},
	{ID: "3", Name: "HDD", Price: 1004.0},
	{ID: "4", Name: "SSD", Price: 1005.0},
	{ID: "5", Name: "GPU", Price: 1006.0},
}
