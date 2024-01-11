package application

import (
	"fmt"
	"net/http"

	"supermarket/internal/handler"
	"supermarket/internal/repository"
	"supermarket/internal/service"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	port string
}

func NewServer(port string) *Server {
	return &Server{port: port}
}

func (s *Server) Run() error {
	rp, _ := repository.NewProductRepository()
	sv := service.NewMovieDefault(rp)
	hd := handler.NewDefaultProducts(sv)

	router := chi.NewRouter()

	router.Get("/ping", handler.Ping)
	router.Route("/products", func(r chi.Router) {
		r.Get("/", hd.GetAllProducts())
		r.Get("/{id}", hd.GetProductById())
		r.Get("/search", hd.GetProductsFiltered())
		r.Post("/", hd.AddProduct())
		r.Put("/{id}", hd.UpdateOrCreateProduct())
		r.Patch("/{id}", hd.PartialProductUpdate())
		r.Delete("/{id}", hd.DeleteProduct())
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%s", s.port), router); err != nil {
		return err
	}
	return nil
}
