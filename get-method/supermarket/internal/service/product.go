package service

import (
	"errors"
	"supermarket/internal"
	"time"
)

type ProductDefault struct {
	repo internal.ProductRepository
}

func NewProductDefault(pdb internal.ProductRepository) *ProductDefault {
	return &ProductDefault{repo: pdb}
}

func (pd *ProductDefault) CheckUniqueCode(code string) (bool, error) {
	_, err := pd.repo.GetByCode(code)
	if err != nil {
		if errors.As(err, &internal.ProductNotFoundError{}) {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (pd *ProductDefault) Save(product internal.Product) internal.Product {
	return pd.repo.Save(product)
}

func (pd *ProductDefault) GetAll() (map[int]internal.Product, error) {
	return pd.repo.GetAll()
}

func (pd *ProductDefault) GetById(id int) (internal.Product, error) {
	return pd.repo.GetById(id)
}

func (pd *ProductDefault) GetByGreaterPrice(price float64) ([]internal.Product, error) {
	return pd.repo.GetByGreaterPrice(price)
}

func (pd *ProductDefault) UpdateOrCreate(product internal.Product) (internal.Product, error) {
	if err := product.Validate(); err != nil {
		return internal.Product{}, err
	}
	return pd.repo.UpdateOrCreate(product)
}

func (pd *ProductDefault) PartialUpdate(id int, product internal.Product) (internal.Product, error) {
	dbProduct, err := pd.repo.GetById(id)
	if err != nil {
		return internal.Product{}, internal.NewProductNotFoundError()
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
		p, err := pd.repo.GetByCode(product.Code)
		if err == nil && p.Id != id {
			return internal.Product{}, internal.NewInvalidProductError("code is not unique")
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
			return internal.Product{}, internal.NewInvalidProductError("expiration")
		}
	}
	return pd.repo.PartialUpdate(id, product)
}

func (pd *ProductDefault) Delete(id int) error {
	return pd.repo.Delete(id)
}

func (pd *ProductDefault) GetTotalPrice(productIds []int) (float64, error) {
	var products []internal.Product
	if len(productIds) == 0 {
		productsMap, err := pd.repo.GetAll()
		if err != nil {
			return 0, err
		}

		for _, product := range productsMap {
			products = append(products, product)
		}
	} else {
		for _, id := range productIds {
			product, err := pd.repo.GetById(id)
			if err != nil {
				return 0, err
			}
			products = append(products, product)
		}
	}

	totalPrice := 0.0
	for _, product := range products {
		totalPrice += product.Price
	}

	return totalPrice, nil
}
