package mux

import (
	"github.com/freekieb7/mud/mux/middleware"
	"net/http"
)

func NewMethodNotAllowedRoute() Route {
	return NewRoute(
		"ANY",
		"/",
		http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
			http.Error(response, "405 method not allowed", http.StatusMethodNotAllowed)
		}),
		make([]middleware.Middleware, 0),
	)
}
