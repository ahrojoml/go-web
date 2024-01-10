package api

type Product struct {
	Id          int     `json:"id,omitempty"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	Code        string  `json:"code_value"`
	IsPublished bool    `json:"is_published,omitempty"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}
