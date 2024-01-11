package internal

type ProductRepository interface {
	Start() (int, error)
	GetAll() (map[int]Product, error)
	GetById(id int) (Product, error)
	Save(product Product) Product
	GetByGreaterPrice(price float64) ([]Product, error)
	GetByCode(code string) (*Product, error)
}
