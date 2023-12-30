package mux

import (
	"net/http"
)

type HandleFunc func(response http.ResponseWriter, request *http.Request)

type router struct {
	tree RouteTree
}

type Router interface {
	Get(path string, handleFunc HandleFunc)
	ServeHTTP(response http.ResponseWriter, request *http.Request)
}

func NewRouter() Router {
	return &router{
		tree: RouteTree{
			root: &RouteNode{
				regex:     "",
				routes:    make([]Route, 0),
				subRoutes: make([]RouteNode, 0),
			},
		},
	}
}

func (router *router) Get(path string, fn HandleFunc) {
	route := NewRoute(http.MethodGet, path, http.HandlerFunc(fn))
	router.tree.Insert(route)
}

func (router *router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	router.tree.Search(request).Handler().ServeHTTP(response, request)
}
