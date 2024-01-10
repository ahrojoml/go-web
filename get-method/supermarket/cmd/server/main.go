package main

import "supermarket/cmd/server/handler"

func main() {
	server := handler.NewServer("8080")
	if err := server.Run(); err != nil {
		panic(err)
	}
}
