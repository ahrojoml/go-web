package internal

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
