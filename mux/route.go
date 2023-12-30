package mux

import (
	"github.com/freekieb7/mud/mux/middleware"
	"log"
	"net/http"
	"regexp"
	"slices"
	"strings"
)

var (
	DynamicArgumentRegex = regexp.MustCompile(`^{.*.}$`) // Starts with { and ends with }
)

type route struct {
	method      string
	path        string
	handler     http.Handler
	middlewares []middleware.Middleware
	params      []string
}

type Route interface {
	Method() string
	Path() string
	Handler() http.Handler
	Middlewares() []middleware.Middleware
}

func NewRoute(method string, path string, handler http.Handler, middlewares ...middleware.Middleware) Route {
	// STEP 1: Cleanup path
	path = strings.Replace(path, "//", "/", -1)

	// STEP 2: Validate path
	pathSlice := strings.Split(path, "/")

	var uniqueEntries []string
	for _, pathEntry := range pathSlice {
		if false == DynamicArgumentRegex.MatchString(pathEntry) {
			continue
		}

		if slices.Contains(uniqueEntries, pathEntry) {
			log.Fatalf("route: duplicate path param found `%s`", pathEntry)
		}

		uniqueEntries = append(uniqueEntries, pathEntry)
	}

	// STEP 3: Return route
	return &route{
		method:      method,
		path:        path,
		middlewares: middlewares,
		handler:     handler,
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
