package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type ProductResponse struct {
	Message string   `json:"message"`
	Data    *Product `json:"data"`
	Error   bool     `json:"error"`
}

type ProductsController struct {
	Products map[int]Product
	LastID   int
}

func (pc *ProductsController) AddProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var product Product
		if err := json.NewDecoder(req.Body).Decode(&product); err != nil {
			body := ProductResponse{
				Message: "could not decode body",
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(body)
		}

		if err := product.validate(); err != nil {
			body := ProductResponse{
				Message: fmt.Sprintf("field is missing body"),
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(body)
			return
		}

		for _, p := range pc.Products {
			if p.Code == product.Code {
				body := ProductResponse{
					Message: "product already exists",
					Error:   true,
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(body)
				return
			}
		}

		if _, err := time.Parse("01/02/2006", product.Expiration); err != nil {
			body := ProductResponse{
				Message: "could not parse expiration date",
				Error:   true,
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(body)
			return
		}

		pc.LastID++
		product.Id = pc.LastID
		pc.Products[product.Id] = product

		body := ProductResponse{
			Message: "success",
			Data:    &product,
			Error:   false,
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(body)
	}
}

func (pc *ProductsController) GetAllProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(pc.Products)
	}
}

func (pc *ProductsController) GetProductById() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(req, "id"))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		for _, product := range pc.Products {
			if product.Id == id {
				w.Header().Set("Content-type", "application/json")
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(product)
			}
		}

		w.WriteHeader(http.StatusNotFound)
	}
}

func (pc *ProductsController) GetProductsFiltered() http.HandlerFunc {
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

		okProducts := []Product{}

		for _, product := range pc.Products {
			if product.Price > price {
				okProducts = append(okProducts, product)
			}
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(okProducts)
	}
}
