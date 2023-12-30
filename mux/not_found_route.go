package mux

import (
	"net/http"
)

func NewNotFoundRoute() Route {
	return NewRoute(
		"ANY",
		"/",
		http.NotFoundHandler(),
	)
}
