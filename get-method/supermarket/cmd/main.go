package main

import "supermarket/internal/application"

func main() {
	server := application.NewServer("8080")
	if err := server.Run(); err != nil {
		panic(err)
	}
}
