package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"supermarket/api"

	"github.com/go-chi/chi/v5"
)

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

	var readProducts []api.Product
	if err := json.Unmarshal(jsonData, &readProducts); err != nil {
		fmt.Println(err)
		return
	}

	var products map[int]api.Product = map[int]api.Product{}
	var nextID int
	for _, product := range readProducts {
		products[product.Id] = product
		nextID = max(nextID, product.Id)
	}

	productController := api.ProductsController{
		Products: products,
		NextID:   nextID + 1,
	}

	router := chi.NewRouter()

	router.Get("/ping", api.Ping)
	router.Route("/products", func(r chi.Router) {
		r.Get("/", productController.GetAllProducts())
		r.Get("/{id}", productController.GetProductById())
		r.Get("/search", productController.GetProductsFiltered())
		r.Post("/", productController.AddProduct())
	})

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}

}
