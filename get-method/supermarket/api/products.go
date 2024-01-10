package api

import (
	"encoding/json"
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
	NextID   int
}

func (pc *ProductsController) AddProduct() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var product Product
		if err := json.NewDecoder(req.Body).Decode(&product); err != nil {
			body := ProductResponse{
				Message: err.Error(),
				Error:   true,
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(body)
		}

		for _, p := range pc.Products {
			if p.Code == product.Code {
				body := ProductResponse{
					Message: "product already exists",
					Error:   true,
				}
				w.WriteHeader(http.StatusBadRequest)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(body)
				return
			}
		}

		if _, err := time.Parse("01/02/2006", product.Expiration); err != nil {
			body := ProductResponse{
				Message: err.Error(),
				Error:   true,
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(body)
			return
		}

		product.Id = pc.NextID
		pc.Products[product.Id] = product
		pc.NextID++

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
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-type", "application/json")
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
				w.WriteHeader(http.StatusOK)
				w.Header().Add("Content-type", "application/json")
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
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode("priceGT value was not set")
			return
		}

		price, err := strconv.ParseFloat(param, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		okProducts := []Product{}

		for _, product := range pc.Products {
			if product.Price > price {
				okProducts = append(okProducts, product)
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-type", "application/json")
		json.NewEncoder(w).Encode(okProducts)
	}
}
