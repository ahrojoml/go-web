package internal

type ProductService interface {
	CheckUniqueCode(code string) (bool, error)
	Save(product Product) Product
	GetAll() (map[int]Product, error)
	GetById(id int) (Product, error)
	GetByGreaterPrice(price float64) ([]Product, error)
}
