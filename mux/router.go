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
	Group(path string, fn func(router Router))
	Use(middleware middleware.Middleware)
	ServeHTTP(response http.ResponseWriter, request *http.Request)
	Routes() []Route
	Middlewares() []middleware.Middleware
}

func NewRouter() Router {
	return &router{
		tree: NewRouteTree(),
	}
}

func (router *router) add(method string, path string, fn HandleFunc, middlewares []middleware.Middleware) {
	if path[0] != '/' {
		path = "/" + path
	}

	newRoute := NewRoute(method, path, http.HandlerFunc(fn), middlewares)
	router.tree.Insert(newRoute)
}

func (router *router) Get(path string, fn HandleFunc, middlewares ...middleware.Middleware) {
	router.add(http.MethodGet, path, fn, middlewares)
}

func (router *router) Post(path string, fn HandleFunc, middlewares ...middleware.Middleware) {
	router.add(http.MethodPost, path, fn, middlewares)
}

func (router *router) Group(path string, fn func(router Router)) {
	subRouter := NewRouter()
	fn(subRouter)

	router.merge(path, subRouter)
}

func (router *router) Use(middleware middleware.Middleware) {
	router.middlewares = append(router.middlewares, middleware)
}

func (router *router) Routes() []Route {
	return router.tree.Routes()
}

func (router *router) Middlewares() []middleware.Middleware {
	return router.middlewares
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

func (router *router) merge(groupPath string, subRouter Router) {
	for _, newRoute := range subRouter.Routes() {
		routePath := newRoute.Path()

		if routePath[0] != '/' {
			routePath = "/" + routePath
		}

		router.add(
			newRoute.Method(),
			groupPath+routePath,
			newRoute.Handler().ServeHTTP,
			append(subRouter.Middlewares(), newRoute.Middlewares()...),
		)
	}
}
