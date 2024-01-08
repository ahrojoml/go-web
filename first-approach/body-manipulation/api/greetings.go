package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Person struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func Greetings(w http.ResponseWriter, r *http.Request) {
	var person Person
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic("aah")
	}

	fmt.Println(body)
	if err := json.Unmarshal(body, &person); err != nil {
		panic("could not unmarshal")
	}

	fmt.Fprint(w, fmt.Sprintf("hello %s %s", person.FirstName, person.LastName))
}
