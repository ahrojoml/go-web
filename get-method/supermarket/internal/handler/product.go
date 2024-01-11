package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"supermarket/internal"

	"supermarket/platform/web/response"

	"github.com/go-chi/chi/v5"
)

type DefaultProducts struct {
	ps internal.ProductService
}

func NewDefaultProducts(ps internal.ProductService) *DefaultProducts {
	return &DefaultProducts{ps: ps}
}

func (pc *DefaultProducts) AddProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var product internal.Product
		if err := json.NewDecoder(req.Body).Decode(&product); err != nil {
			response.Error(w, http.StatusBadRequest, "could not decode body")
			return
		}

		if err := product.Validate(); err != nil {
			response.Error(w, http.StatusBadRequest, "field is missing")
			return
		}

		productExists, err := pc.ps.CheckUniqueCode(product.Code)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "error retrieving product by code")
			return
		}

		if !productExists {
			response.Error(w, http.StatusConflict, "product already exists")
			return
		}

		product = pc.ps.Save(product)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
	}
}

func (pc *DefaultProducts) GetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		products, err := pc.ps.GetAll()
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "error retrieving products")
			return
		}
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(products)
	}
}

func (pc *DefaultProducts) GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(req, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, "error parsing id")
			return
		}

		product, err := pc.ps.GetById(id)
		if err != nil {
			response.Error(w, http.StatusNotFound, "product not found")
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(product)

	}
}

func (pc *DefaultProducts) GetProductsFiltered() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		param := req.URL.Query().Get("priceGT")
		if param == "" {
			response.Error(w, http.StatusBadRequest, "priceGT value was not set")
			return
		}

		price, err := strconv.ParseFloat(param, 64)
		if err != nil {
			response.Error(w, http.StatusBadRequest, "error parsing priceGT value")
			return
		}

		products, err := pc.ps.GetByGreaterPrice(price)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "error retrieving products")
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(products)
	}
}

func (pc *DefaultProducts) UpdateOrCreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var product internal.Product
		if err := json.NewDecoder(req.Body).Decode(&product); err != nil {
			response.Error(w, http.StatusBadRequest, "could not decode body")
			return
		}

		updatedProduct, err := pc.ps.UpdateOrCreate(product)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "error updating or creating product")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedProduct)
	}
}

func (pc *DefaultProducts) PartialProductUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(req, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, "error parsing id")
			return
		}

		var product internal.Product
		if err := json.NewDecoder(req.Body).Decode(&product); err != nil {
			response.Error(w, http.StatusBadRequest, "could not decode body")
			return
		}

		product.Id = id

		updatedProduct, err := pc.ps.PartialUpdate(id, product)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "error updating product")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(updatedProduct)
	}
}

func (pc *DefaultProducts) DeleteProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(req, "id"))
		if err != nil {
			response.Error(w, http.StatusBadRequest, "error parsing id")
			return
		}

		err = pc.ps.Delete(id)
		if err != nil {
			if errors.As(err, &internal.ProductNotFoundError{}) {
				response.Error(w, http.StatusNotFound, "product not found")
				return
			}
			response.Error(w, http.StatusInternalServerError, "error deleting product")
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
