package application

import (
	"fmt"
	"net/http"

	"supermarket/internal/handler"
	"supermarket/internal/repository"
	"supermarket/internal/service"

	"supermarket/platform/web/middleware"

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
	sv := service.NewProductDefault(rp)
	hd := handler.NewDefaultProducts(sv)

	router := chi.NewRouter()

	router.Use(middleware.ResponseData)
	router.Get("/ping", handler.Ping)
	router.Route("/products", func(r chi.Router) {
		r.Get("/", hd.GetAllProducts())
		r.Get("/{id}", hd.GetProductById())
		r.Get("/search", hd.GetProductsFiltered())
		r.With(middleware.Auth).Post("/", hd.AddProduct())
		r.With(middleware.Auth).Put("/", hd.UpdateOrCreateProduct())
		r.With(middleware.Auth).Patch("/{id}", hd.PartialProductUpdate())
		r.With(middleware.Auth).Delete("/{id}", hd.DeleteProduct())
		r.Get("/consumer_price", hd.GetCartPrice())
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%s", s.port), router); err != nil {
		return err
	}
	return nil
}
