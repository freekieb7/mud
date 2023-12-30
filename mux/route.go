package mux

import (
	"net/http"
	"regexp"
)

var (
	DynamicArgumentRegex = regexp.MustCompile(`^{.*.}$`) // Starts with { and ends with }
)

type route struct {
	method  string
	path    string
	handler http.Handler
}

type Route interface {
	Method() string
	Path() string
	Handler() http.Handler
}

func NewRoute(method string, path string, handler http.Handler) Route {
	return &route{
		method:  method,
		path:    path,
		handler: handler,
	}
}

func (route *route) Method() string {
	return route.method
}

func (route *route) Path() string {
	return route.path
}

func (route *route) Handler() http.Handler {
	return route.handler
}
