package mux

import (
	"github.com/freekieb7/mux/middleware"
	"net/http"
	"regexp"
)

var (
	DynamicArgumentRegex = regexp.MustCompile(`^{.*.}$`) // Starts with { and ends with }
)

type route struct {
	method      string
	path        string
	handler     http.Handler
	middlewares []middleware.Middleware
}

type Route interface {
	Method() string
	Path() string
	Handler() http.Handler
	Middlewares() []middleware.Middleware
}

func NewRoute(method string, path string, handler http.Handler, middlewares []middleware.Middleware) Route {
	return &route{
		method:      method,
		path:        path,
		handler:     handler,
		middlewares: middlewares,
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

func (route *route) Middlewares() []middleware.Middleware {
	return route.middlewares
}
