package product

import (
	"encoding/json"
	"io"
	"os"
)

// This will be the db
var products map[int]Product = map[int]Product{}
var lastID int = 0

func loadDB() (int, error) {
	file, err := os.Open("./get-method/supermarket/docs/db/products.json")
	if err != nil {
		return 0, err
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	var readProducts []Product
	if err := json.Unmarshal(jsonData, &readProducts); err != nil {
		return 0, err
	}

	for _, product := range readProducts {
		products[product.Id] = product
		lastID = max(lastID, product.Id)
	}

	return lastID, nil
}

func getAll() (map[int]Product, error) {
	return products, nil
}

func getById(id int) (Product, error) {
	for _, product := range products {
		if product.Id == id {
			return product, nil
		}
	}
	return Product{}, NewProductNotFoundError()
}

func save(product Product) Product {
	lastID++
	product.Id = lastID
	products[lastID] = product
	return product
}

func getByGreaterPrice(price float64) ([]Product, error) {
	okProducts := []Product{}

	for _, product := range products {
		if product.Price > price {
			okProducts = append(okProducts, product)
		}
	}
	return okProducts, nil
}

func getByCode(code string) (*Product, error) {
	for _, product := range products {
		if product.Code == code {
			return &product, nil
		}
	}
	return nil, NewProductNotFoundError()
}
