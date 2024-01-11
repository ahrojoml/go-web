package repository

import (
	"encoding/json"
	"io"
	"os"
	"supermarket/internal"
	"time"
)

type ProductDB struct {
	Products map[int]internal.Product
	LastID   int
}

func NewProductRepository() (*ProductDB, error) {
	pdb := &ProductDB{
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

func (pdb *ProductDB) Start() (int, error) {
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

func (pdb *ProductDB) GetAll() (map[int]internal.Product, error) {
	return pdb.Products, nil
}

func (pdb *ProductDB) GetById(id int) (internal.Product, error) {
	for _, product := range pdb.Products {
		if product.Id == id {
			return product, nil
		}
	}
	return internal.Product{}, internal.NewProductNotFoundError()
}

func (pdb *ProductDB) Save(product internal.Product) internal.Product {
	pdb.LastID++
	product.Id = pdb.LastID
	pdb.Products[pdb.LastID] = product
	return product
}

func (pdb *ProductDB) GetByGreaterPrice(price float64) ([]internal.Product, error) {
	okProducts := []internal.Product{}

	for _, product := range pdb.Products {
		if product.Price > price {
			okProducts = append(okProducts, product)
		}
	}
	return okProducts, nil
}

func (pdb *ProductDB) GetByCode(code string) (*internal.Product, error) {
	for _, product := range pdb.Products {
		if product.Code == code {
			return &product, nil
		}
	}
	return nil, internal.NewProductNotFoundError()
}

func (pdb *ProductDB) UpdateOrCreate(product internal.Product) (*internal.Product, error) {
	if err := product.Validate(); err != nil {
		return nil, err
	}

	if product.Id == 0 {
		p, err := pdb.GetByCode(product.Code)
		if err == nil && p.Id != product.Id {
			return nil, internal.NewInvalidProductError("code is not unique")
		}
		pdb.Save(product)
		return &product, nil
	}

	pdb.Products[product.Id] = product
	return &product, nil
}

func (pdb *ProductDB) PartialUpdate(id int, product internal.Product) (*internal.Product, error) {
	dbProduct, ok := pdb.Products[id]
	if !ok {
		return nil, internal.NewProductNotFoundError()
	}

	if product.Name == "" {
		product.Name = dbProduct.Name
	}

	if product.Quantity == 0 {
		product.Quantity = dbProduct.Quantity
	}

	if product.Code == "" {
		product.Code = dbProduct.Code
	} else {
		p, err := pdb.GetByCode(product.Code)
		if err == nil && p.Id != id {
			return nil, internal.NewInvalidProductError("code is not unique")
		}
	}

	if product.Price == 0 {
		product.Price = dbProduct.Price
	}

	if product.IsPublished == false {
		product.IsPublished = dbProduct.IsPublished
	}

	if product.Expiration == "" {
		product.Expiration = dbProduct.Expiration
	} else {
		_, err := time.Parse("02/01/2006", product.Expiration)
		if err != nil {
			return nil, internal.NewInvalidProductError("expiration")
		}
	}

	pdb.Products[id] = product

	return &product, nil

}

func (pdb *ProductDB) Delete(id int) error {
	_, ok := pdb.Products[id]
	if !ok {
		return internal.NewProductNotFoundError()
	}

	delete(pdb.Products, id)
	return nil
}
