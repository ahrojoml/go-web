package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
)

type Product struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	Code        string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

func GetAllProducts(w http.ResponseWriter, req *http.Request) {
	response, err := json.Marshal(products)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "application/json")
	fmt.Fprint(w, string(response))
}

func GetProductById(w http.ResponseWriter, req *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, product := range products {
		if product.Id == id {
			response, err := json.Marshal(product)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprint(w, err.Error())
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Header().Add("Content-type", "application/json")
			fmt.Fprint(w, string(response))
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func GetProductsFiltered(w http.ResponseWriter, req *http.Request) {
	param := req.URL.Query().Get("priceGT")
	if param == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "priceGT value was not set")
		return
	}

	price, err := strconv.ParseFloat(param, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err.Error())
		return
	}

	okProducts := []Product{}

	for _, product := range products {
		if product.Price > price {
			okProducts = append(okProducts, product)
		}
	}

	response, err := json.Marshal(okProducts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-type", "application/json")
	fmt.Fprint(w, string(response))
}

func Ping(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "pong")
}

var products []Product

func main() {
	file, err := os.Open("./get-method/supermarket/products.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := json.Unmarshal(jsonData, &products); err != nil {
		fmt.Println(err)
		return
	}

	router := chi.NewRouter()

	router.Get("/ping", Ping)
	router.Get("/products", GetAllProducts)
	router.Get("/products/{id}", GetProductById)
	router.Get("/products/search", GetProductsFiltered)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}

}
