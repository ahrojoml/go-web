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

func (p Product) validate() error {
	if p.Name == "" {
		return NewInvalidProductError("name")
	}
	if p.Quantity == 0 {
		return NewInvalidProductError("quantity")
	}
	if p.Code == "" {
		return NewInvalidProductError("code")
	}
	if p.Expiration == "" {
		return NewInvalidProductError("expiration")
	}
	if p.Price == 0 {
		return NewInvalidProductError("price")
	}
	return nil
}
