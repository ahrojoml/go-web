package product

import "errors"

func LoadDB() (int, error) {
	return loadDB()
}

func CheckUniqueCode(code string) (bool, error) {
	_, err := getByCode(code)
	if err != nil {
		if errors.As(err, &ProductNotFoundError{}) {
			return true, nil
		}
		return false, nil
	}
	return false, nil
}

func Save(product Product) Product {
	return save(product)
}

func GetAll() (map[int]Product, error) {
	return getAll()
}

func GetById(id int) (Product, error) {
	return getById(id)
}

func GetByGreaterPrice(price float64) ([]Product, error) {
	return getByGreaterPrice(price)
}
