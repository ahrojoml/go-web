package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"supermarket/internal"
	"supermarket/internal/handler"
	"supermarket/internal/repository"
	"supermarket/internal/service"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestGetAllProducts(t *testing.T) {
	dbData := map[int]internal.Product{
		1: {Id: 1, Name: "p1", Quantity: 1, Code: "c1", Price: 1, IsPublished: true},
		2: {Id: 2, Name: "p2", Quantity: 2, Code: "c2", Price: 2, IsPublished: true},
	}
	db := repository.ProductMapDB{Products: dbData, LastID: 2}
	sv := service.NewProductDefault(&db)
	hd := handler.NewDefaultProducts(sv)

	expectCode := http.StatusOK
	expectBody, _ := json.Marshal(dbData)
	expectHeader := http.Header{"Content-Type": []string{"application/json"}}

	req := httptest.NewRequest("GET", "/products", nil)
	res := httptest.NewRecorder()
	hd.GetAllProducts()(res, req)

	require.Equal(t, expectCode, res.Code)
	require.JSONEq(t, string(expectBody), res.Body.String())
	require.Equal(t, expectHeader, res.Header())
}

func TestGetProduct(t *testing.T) {
	dbData := map[int]internal.Product{
		1: {Id: 1, Name: "p1", Quantity: 1, Code: "c1", Price: 1, IsPublished: true},
		2: {Id: 2, Name: "p2", Quantity: 2, Code: "c2", Price: 2, IsPublished: true},
	}
	db := repository.ProductMapDB{Products: dbData, LastID: 2}
	sv := service.NewProductDefault(&db)
	hd := handler.NewDefaultProducts(sv)

	expectCode := http.StatusOK
	expectBody, _ := json.Marshal(dbData[1])
	expectHeader := http.Header{"Content-Type": []string{"application/json"}}

	req := httptest.NewRequest("GET", "/products/1/", nil)

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	res := httptest.NewRecorder()
	hd.GetProductById()(res, req)

	require.Equal(t, expectCode, res.Code)
	require.JSONEq(t, string(expectBody), res.Body.String())
	require.Equal(t, expectHeader, res.Header())
}

func TestAddProduct(t *testing.T) {
	dbData := map[int]internal.Product{
		1: {Id: 1, Name: "p1", Quantity: 1, Code: "c1", Price: 1, IsPublished: true},
		2: {Id: 2, Name: "p2", Quantity: 2, Code: "c2", Price: 2, IsPublished: true},
	}
	db := repository.ProductMapDB{Products: dbData, LastID: 2}
	sv := service.NewProductDefault(&db)
	hd := handler.NewDefaultProducts(sv)

	newProd := internal.Product{
		Id: 3, Name: "p3", Quantity: 3, Code: "c3", Price: 3, IsPublished: false, Expiration: "01/02/2065",
	}

	expectCode := http.StatusCreated
	expectBody, _ := json.Marshal(newProd)
	expectHeader := http.Header{"Content-Type": []string{"application/json"}}

	req := httptest.NewRequest("POST", "/products", strings.NewReader(string(expectBody)))
	res := httptest.NewRecorder()
	hd.AddProduct()(res, req)

	require.Equal(t, expectCode, res.Code)
	require.JSONEq(t, string(expectBody), res.Body.String())
	require.Equal(t, expectHeader, res.Header())
	require.Equal(t, 3, len(db.Products))
}

func TestDeleteProduct(t *testing.T) {
	dbData := map[int]internal.Product{
		1: {Id: 1, Name: "p1", Quantity: 1, Code: "c1", Price: 1, IsPublished: true},
		2: {Id: 2, Name: "p2", Quantity: 2, Code: "c2", Price: 2, IsPublished: true},
	}
	db := repository.ProductMapDB{Products: dbData, LastID: 2}
	sv := service.NewProductDefault(&db)
	hd := handler.NewDefaultProducts(sv)

	expectCode := http.StatusOK
	expectHeader := http.Header{"Content-Type": []string{"application/json"}}

	req := httptest.NewRequest("DELETE", "/products/1/", nil)

	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", "1")

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))

	res := httptest.NewRecorder()
	hd.DeleteProduct()(res, req)

	require.Equal(t, expectCode, res.Code)
	require.Equal(t, expectHeader, res.Header())
	require.Equal(t, 1, len(db.Products))

	_, ok := db.Products[1]
	require.False(t, ok)
}
