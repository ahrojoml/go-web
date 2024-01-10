package application

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"supermarket/api"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	port string
}

func NewServer(port string) *Server {
	return &Server{port: port}
}

func (s *Server) Run() error {
	file, err := os.Open("./get-method/supermarket/products.json")
	if err != nil {
		return err
	}

	jsonData, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var readProducts []api.Product
	if err := json.Unmarshal(jsonData, &readProducts); err != nil {
		return err
	}

	var products map[int]api.Product = map[int]api.Product{}
	var lastID int
	for _, product := range readProducts {
		products[product.Id] = product
		lastID = max(lastID, product.Id)
	}

	productController := api.ProductsController{
		Products: products,
		LastID:   lastID,
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
		return err
	}
	return nil
}
