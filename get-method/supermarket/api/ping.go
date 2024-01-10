package api

import (
	"encoding/json"
	"net/http"
)

func Ping(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("pong")
}
