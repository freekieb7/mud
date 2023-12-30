package mux

import (
	"net/http"
)

func NewInternalServerErrorRoute() Route {
	return NewRoute(
		"ANY",
		"",
		http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			http.Error(response, "500 internal server error", http.StatusInternalServerError)
		}),
	)
}
