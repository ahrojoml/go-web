package main

import (
	"fmt"
	"net/http"
)

func Ping(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "pong")
}

func main() {
	http.HandleFunc("/ping", Ping)

	http.ListenAndServe(":8081", nil)
}
