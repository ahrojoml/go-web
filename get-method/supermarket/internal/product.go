package internal

import "time"

type ProductRepository interface {
	Start() (int, error)
	GetAll() (map[int]Product, error)
	GetById(id int) (Product, error)
	Save(product Product) Product
	GetByGreaterPrice(price float64) ([]Product, error)
	GetByCode(code string) (*Product, error)
	UpdateOrCreate(product Product) (Product, error)
	PartialUpdate(id int, product Product) (Product, error)
	Delete(id int) error
}

type ProductDBRepository interface {
	GetById(id int) (Product, error)
	Create(product *Product) error
	Update(product *Product) error
	Delete(id int) error
}

type ProductService interface {
	CheckUniqueCode(code string) (bool, error)
	Save(product Product) Product
	GetAll() (map[int]Product, error)
	GetById(id int) (Product, error)
	GetByGreaterPrice(price float64) ([]Product, error)
	UpdateOrCreate(product Product) (Product, error)
	PartialUpdate(id int, product Product) (Product, error)
	Delete(id int) error
	GetTotalPrice(productIds []int) (float64, error)
}

type InvalidProductError struct {
	Field string
}

func (e InvalidProductError) Error() string {
	return "invalid product"
}

func NewInvalidProductError(field string) error {
	return InvalidProductError{Field: field}
}

type ProductNotFoundError struct{}

func (e ProductNotFoundError) Error() string {
	return "product not found"
}

func NewProductNotFoundError() error {
	return ProductNotFoundError{}
}

type ProductAlreadyExistsError struct{}

func (e ProductAlreadyExistsError) Error() string {
	return "product already exists"
}

func NewProductAlreadyExistsError() error {
	return ProductAlreadyExistsError{}
}

type Product struct {
	Id          int     `json:"id,omitempty"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	Code        string  `json:"code_value"`
	IsPublished bool    `json:"is_published,omitempty"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func (p Product) Validate() error {
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
	if p.Price <= 0 {
		return NewInvalidProductError("price")
	}
	if _, err := time.Parse("02/01/2006", p.Expiration); err != nil {
		return NewInvalidProductError("date")
	}
	return nil
}
