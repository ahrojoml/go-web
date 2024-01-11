package main

import (
	"os"
	"supermarket/internal/application"
)

func main() {
	os.Setenv("PRODUCT_KEY", "123")
	os.Setenv("DB_FILE_PATH", "./get-method/supermarket/docs/db/products.json")

	server := application.NewServer("8080")
	if err := server.Run(); err != nil {
		panic(err)
	}
}
