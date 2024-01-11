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

	file, err := os.Open("./get-method/supermarket/docs/db/products.json")
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

func (s *Storage) Save(map[int]internal.Product) error {
	return nil
}
