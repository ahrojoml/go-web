package repository

import (
	"encoding/json"
	"io"
	"os"
	"supermarket/internal"
)

type Storage struct{}

func (s *Storage) GetAll() (map[int]internal.Product, int, error) {
	var products map[int]internal.Product = map[int]internal.Product{}

	file, err := os.Open(os.Getenv("DB_PATH"))
	if err != nil {
		return map[int]internal.Product{}, 0, err
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return map[int]internal.Product{}, 0, err
	}

	var readProducts []internal.Product
	if err := json.Unmarshal(jsonData, &readProducts); err != nil {
		return map[int]internal.Product{}, 0, err
	}

	var lastId int = 0
	for _, product := range readProducts {
		products[product.Id] = product
		lastId = max(lastId, product.Id)
	}

	return map[int]internal.Product{}, lastId, nil
}

func (s *Storage) Save(products map[int]internal.Product) error {
	file, err := os.Create(os.Getenv("DB_PATH"))
	if err != nil {
		return err
	}

	var productsArray []internal.Product = make([]internal.Product, 0, len(products))
	for _, product := range products {
		productsArray = append(productsArray, product)
	}

	jsonData, err := json.Marshal(productsArray)
	if err != nil {
		return err
	}

	if _, err := file.Write(jsonData); err != nil {
		return err
	}

	return nil
}
