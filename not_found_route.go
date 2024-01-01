package mux

import (
	"github.com/freekieb7/mux/middleware"
	"net/http"
)

func NewNotFoundRoute() Route {
	return NewRoute(
		"ANY",
		"",
		http.NotFoundHandler(),
		make([]middleware.Middleware, 0),
	)
}
