package middleware

import (
	"fmt"
	"net/http"
	"time"
)

func ResponseData(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r)

		fmt.Printf("Method: %s\n", r.Method)
		fmt.Printf("Time: %s\n", time.Now())
		fmt.Printf("Path: %s%s\n", r.Host, r.URL.Path)
		fmt.Printf("Request Size: %d\n\n", r.ContentLength)

	})
}
