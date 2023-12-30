package mux

import (
	"net/http"
)

func NewMethodNotAllowedRoute() Route {
	return NewRoute(
		"ANY",
		"/",
		http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			http.Error(response, "405 method not allowed", http.StatusMethodNotAllowed)
		}),
	)
}
