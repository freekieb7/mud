package mux

import (
	"github.com/freekieb7/mud/mux/middleware"
	"net/http"
)

type HandleFunc func(response http.ResponseWriter, request *http.Request)

type router struct {
	tree        RouteTree
	middlewares []middleware.Middleware
}

type Router interface {
	Get(path string, handleFunc HandleFunc, middlewares ...middleware.Middleware)
	Post(path string, handleFunc HandleFunc, middlewares ...middleware.Middleware)
	Use(middleware middleware.Middleware)
	ServeHTTP(response http.ResponseWriter, request *http.Request)
}

func NewRouter() Router {
	return &router{
		tree: RouteTree{
			root: RouteNode{
				regex:     "",
				routes:    []Route{},
				subRoutes: map[string]*RouteNode{},
			},
		},
	}
}

func (router *router) add(method string, path string, fn HandleFunc, middlewares []middleware.Middleware) {
	newRoute := NewRoute(method, path, http.HandlerFunc(fn), middlewares)
	router.tree.Insert(newRoute)
}

func (router *router) Get(path string, fn HandleFunc, middlewares ...middleware.Middleware) {
	router.add(http.MethodGet, path, fn, middlewares)
}

func (router *router) Post(path string, fn HandleFunc, middlewares ...middleware.Middleware) {
	router.add(http.MethodPost, path, fn, middlewares)
}

func (router *router) Use(middleware middleware.Middleware) {
	router.middlewares = append(router.middlewares, middleware)
}

func (router *router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	matchingRoute := router.tree.Search(request)
	handler := matchingRoute.Handler()

	// Route specific middleware
	for _, routeMiddleware := range matchingRoute.Middlewares() {
		handler = routeMiddleware.Process(handler)
	}

	// Go through router middleware
	for _, routerMiddleware := range router.middlewares {
		handler = routerMiddleware.Process(handler)
	}

	handler.ServeHTTP(response, request)
}
