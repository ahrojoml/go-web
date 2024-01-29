package repository

import (
	"encoding/json"
	"io"
	"os"
	"supermarket/internal"
)

type ProductMapDB struct {
	Products map[int]internal.Product
	LastID   int
}

func NewProductRepository() (*ProductMapDB, error) {
	pdb := &ProductMapDB{
		Products: map[int]internal.Product{},
		LastID:   0,
	}
	lastId, err := pdb.Start()
	if err != nil {
		return nil, err
	}
	pdb.LastID = lastId
	return pdb, nil
}

func (pdb *ProductMapDB) Start() (int, error) {
	file, err := os.Open("./get-method/supermarket/docs/db/products.json")
	if err != nil {
		return 0, err
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return 0, err
	}

	var readProducts []internal.Product
	if err := json.Unmarshal(jsonData, &readProducts); err != nil {
		return 0, err
	}

	var lastId int = 0
	for _, product := range readProducts {
		pdb.Products[product.Id] = product
		lastId = max(lastId, product.Id)
	}

	return lastId, nil
}

func (pdb *ProductMapDB) GetAll() (map[int]internal.Product, error) {
	return pdb.Products, nil
}

func (pdb *ProductMapDB) GetById(id int) (internal.Product, error) {
	for _, product := range pdb.Products {
		if product.Id == id {
			return product, nil
		}
	}
	return internal.Product{}, internal.NewProductNotFoundError()
}

func (pdb *ProductMapDB) Save(product internal.Product) internal.Product {
	pdb.LastID++
	product.Id = pdb.LastID
	pdb.Products[pdb.LastID] = product
	return product
}

func (pdb *ProductMapDB) GetByGreaterPrice(price float64) ([]internal.Product, error) {
	okProducts := []internal.Product{}

	for _, product := range pdb.Products {
		if product.Price > price {
			okProducts = append(okProducts, product)
		}
	}
	return okProducts, nil
}

func (pdb *ProductMapDB) GetByCode(code string) (*internal.Product, error) {
	for _, product := range pdb.Products {
		if product.Code == code {
			return &product, nil
		}
	}
	return nil, internal.NewProductNotFoundError()
}

func (pdb *ProductMapDB) UpdateOrCreate(product internal.Product) (internal.Product, error) {
	if product.Id == 0 {
		p, err := pdb.GetByCode(product.Code)
		if err == nil && p.Id != product.Id {
			return internal.Product{}, internal.NewInvalidProductError("code is not unique")
		}
		pdb.Save(product)
		product.Id = pdb.LastID
		return product, nil
	}

	pdb.Products[product.Id] = product
	return product, nil
}

func (pdb *ProductMapDB) PartialUpdate(id int, product internal.Product) (internal.Product, error) {

	pdb.Products[id] = product

	return product, nil

}

func (pdb *ProductMapDB) Delete(id int) error {
	_, ok := pdb.Products[id]
	if !ok {
		return internal.NewProductNotFoundError()
	}

	delete(pdb.Products, id)
	return nil
}
