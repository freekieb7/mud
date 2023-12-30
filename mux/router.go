package mux

import (
	"context"
	"github.com/freekieb7/mud/mux/middleware"
	"net/http"
)

type HandleFunc func(response http.ResponseWriter, request *http.Request)

type router struct {
	basePath    string
	routes      []Route
	matchers    []Matcher
	middlewares []middleware.Middleware
}

type Router interface {
	Routes() []Route
	Middlewares() []middleware.Middleware
	Get(path string, handleFunc HandleFunc, middlewares ...middleware.Middleware)
	Post(path string, handleFunc HandleFunc, middlewares ...middleware.Middleware)
	Group(path string, fn func(router Router))
	Use(middleware middleware.Middleware)
	ServeHTTP(response http.ResponseWriter, request *http.Request)
}

type RouterConfig struct {
	Matchers    []Matcher
	Routes      []Route
	Middlewares []middleware.Middleware
}

func NewRouter(config RouterConfig) Router {
	return &router{
		matchers: config.Matchers,
		routes:   config.Routes,
	}
}

func NewDefaultRouter() Router {
	return NewRouter(
		RouterConfig{
			Matchers: []Matcher{
				NewMethodMatcher(),
				NewPathMatcher(),
			},
			Middlewares: []middleware.Middleware{
				middleware.NewDefaultLoggingMiddleware(),
			},
		},
	)
}

func (router *router) AddRoute(method string, path string, handleFunc http.Handler, middlewares ...middleware.Middleware) {
	newRoute := NewRoute(method, path, handleFunc, middlewares...)
	router.routes = append(router.routes, newRoute)
}

func (router *router) Routes() []Route {
	return router.routes
}

func (router *router) Middlewares() []middleware.Middleware {
	return router.middlewares
}

func (router *router) Get(path string, handleFunc HandleFunc, middlewares ...middleware.Middleware) {
	router.AddRoute(http.MethodGet, path, http.HandlerFunc(handleFunc), middlewares...)
}

func (router *router) Post(path string, handleFunc HandleFunc, middlewares ...middleware.Middleware) {
	router.AddRoute(http.MethodPost, path, http.HandlerFunc(handleFunc), middlewares...)
}

func (router *router) Group(path string, fn func(router Router)) {
	subRouter := NewDefaultRouter()
	fn(subRouter)

	router.merge(path, subRouter)
}

func (router *router) Use(middleware middleware.Middleware) {
	router.middlewares = append(router.middlewares, middleware)
}

func (router *router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// STEP 1: Pick correct route
	var targetRoute Route

	for _, routerRoute := range router.routes {
		if false == router.match(routerRoute, request) {
			continue
		}

		targetRoute = routerRoute
		break
	}

	if targetRoute.Handler == nil {
		targetRoute = NewNotFoundRoute()
	}

	// STEP 2: Add route to request context
	request = request.WithContext(context.WithValue(request.Context(), RouteCtxKey, targetRoute))

	// STEP 3: Run through middleware
	handler := targetRoute.Handler()

	/// Go through route specific middleware
	for _, routeMiddleware := range targetRoute.Middlewares() {
		handler = routeMiddleware.Process(handler)
	}

	/// Go through router middleware
	for _, routerMiddleware := range router.middlewares {
		handler = routerMiddleware.Process(handler)
	}

	// STEP 4: Run route handler
	handler.ServeHTTP(response, request)
}

func (router *router) match(route Route, request *http.Request) bool {
	for _, routerMatcher := range router.matchers {
		if false == routerMatcher.Match(route, request) {
			return false
		}
	}

	return true
}

func (router *router) merge(path string, subRouter Router) {
	for _, routerRoute := range subRouter.Routes() {
		router.AddRoute(
			routerRoute.Method(),
			path+routerRoute.Path(),
			routerRoute.Handler(),
			append(subRouter.Middlewares(), routerRoute.Middlewares()...)...,
		)
	}
}
