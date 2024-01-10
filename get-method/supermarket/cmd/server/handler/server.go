package handler

import (
	"fmt"
	"net/http"
	"supermarket/internal/product"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	port string
}

func NewServer(port string) *Server {
	return &Server{port: port}
}

func (s *Server) Run() error {
	productController := ProductsController{}

	product.LoadDB()

	router := chi.NewRouter()

	router.Get("/ping", Ping)
	router.Route("/products", func(r chi.Router) {
		r.Get("/", productController.GetAllProducts())
		r.Get("/{id}", productController.GetProductById())
		r.Get("/search", productController.GetProductsFiltered())
		r.Post("/", productController.AddProduct())
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%s", s.port), router); err != nil {
		return err
	}
	return nil
}
