package service

import (
	"errors"
	"supermarket/internal"
)

type ProductDefault struct {
	repo internal.ProductRepository
}

func NewMovieDefault(pdb internal.ProductRepository) *ProductDefault {
	return &ProductDefault{repo: pdb}
}

func (pd *ProductDefault) CheckUniqueCode(code string) (bool, error) {
	_, err := pd.repo.GetByCode(code)
	if err != nil {
		if errors.As(err, &internal.ProductNotFoundError{}) {
			return true, nil
		}
		return false, nil
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

func (pd *ProductDefault) UpdateOrCreate(product internal.Product) (*internal.Product, error) {
	return pd.repo.UpdateOrCreate(product)
}

func (pd *ProductDefault) PartialUpdate(id int, product internal.Product) (*internal.Product, error) {
	return pd.repo.PartialUpdate(id, product)
}

func (pd *ProductDefault) Delete(id int) error {
	return pd.repo.Delete(id)
}
