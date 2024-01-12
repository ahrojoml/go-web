package middleware

import (
	"crypto/subtle"
	"net/http"
	"os"
	"supermarket/platform/web/response"
)

func Auth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token := r.Header.Get("Authorization"); subtle.ConstantTimeCompare([]byte(token), []byte(os.Getenv("PRODUCT_KEY"))) != 1 {
			response.Error(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		handler.ServeHTTP(w, r)
	})
}
