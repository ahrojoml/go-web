package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"supermarket/internal"

	"github.com/go-chi/chi/v5"
)

type DefaultProducts struct {
	ps internal.ProductService
}

func NewDefaultProducts(ps internal.ProductService) *DefaultProducts {
	return &DefaultProducts{ps: ps}
}

type ProductResponse struct {
	Message string            `json:"message"`
	Data    *internal.Product `json:"data"`
	Error   bool              `json:"error"`
}

func (pc *DefaultProducts) AddProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var product internal.Product
		if err := json.NewDecoder(req.Body).Decode(&product); err != nil {
			body := ProductResponse{
				Message: "could not decode body",
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(body)
		}

		if err := product.Validate(); err != nil {
			body := ProductResponse{
				Message: fmt.Sprintf("field is missing body"),
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(body)
			return
		}

		productExists, err := pc.ps.CheckUniqueCode(product.Code)
		if err != nil {
			body := ProductResponse{
				Message: "error retrieving product by code",
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(body)
			return
		}

		if !productExists {
			body := ProductResponse{
				Message: "product already exists",
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(body)
			return
		}

		product = pc.ps.Save(product)

		body := ProductResponse{
			Message: "success",
			Data:    &product,
			Error:   false,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(body)
	}
}

func (pc *DefaultProducts) GetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		products, err := pc.ps.GetAll()
		if err != nil {
			body := ProductResponse{
				Message: "error retrieving products",
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(body)
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
			w.WriteHeader(http.StatusNotFound)
			return
		}

		product, err := pc.ps.GetById(id)
		if err != nil {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode("product not found")
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
			json.NewEncoder(w).Encode("priceGT value was not set")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		price, err := strconv.ParseFloat(param, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("error parsing priceGT value")
			return
		}

		products, err := pc.ps.GetByGreaterPrice(price)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode("error retrieving products")
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(products)
	}
}

func (pc *DefaultProducts) UpdateOrCreateProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(req, "id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var product internal.Product
		if err := json.NewDecoder(req.Body).Decode(&product); err != nil {
			body := ProductResponse{
				Message: "could not decode body",
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(body)
		}

		product.Id = id

		updatedProduct, err := pc.ps.UpdateOrCreate(product)
		if err != nil {
			body := ProductResponse{
				Message: "error updating or creating product",
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(body)
			return
		}

		body := ProductResponse{
			Message: "success",
			Data:    updatedProduct,
			Error:   false,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(body)
	}
}
