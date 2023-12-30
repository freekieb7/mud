package mux

import (
	"log"
	"net/http"
	"strings"
)

type RouteTree struct {
	root *RouteNode
}

type RouteNode struct {
	regex     string
	routes    []Route
	subRoutes []RouteNode
}

func (tree *RouteTree) Insert(route Route) {
	pathPieces := strings.Split(route.Path(), "/")

	tree.root.insert(route, pathPieces)
}

func (node *RouteNode) insert(route Route, pathPieces []string) {
	if node.regex != pathPieces[0] {
		log.Fatal("something is not right")
	}

	// Match made in heaven
	if len(pathPieces) == 1 {
		node.routes = append(node.routes, route)
		return
	}

	// Maybe subroute
	for index := range node.subRoutes {
		if node.subRoutes[index].regex == pathPieces[1] {
			node.subRoutes[index].insert(route, pathPieces[1:])
		}
	}

	// No subroute available so create new
	node.subRoutes = append(node.subRoutes, RouteNode{
		regex:     pathPieces[1],
		routes:    make([]Route, 0),
		subRoutes: make([]RouteNode, 0),
	})
	node.subRoutes[len(node.subRoutes)-1].insert(route, pathPieces[1:])
}

func (tree *RouteTree) Search(request *http.Request) Route {
	pathPieces := strings.Split(request.URL.Path, "/")
	return tree.root.search(pathPieces)
}

func (node *RouteNode) search(pathPieces []string) Route {
	if len(pathPieces) == 0 {
		return NewRoute(http.MethodGet, "", http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			http.Error(writer, "500 internal server error", http.StatusInternalServerError)
		}))
	}

	if node.regex != pathPieces[0] {
		return NewRoute(http.MethodGet, "", http.NotFoundHandler())
	}

	// Match made in heaven
	if len(pathPieces) == 1 {
		return node.routes[0]
	}

	for index := range node.subRoutes {
		if node.subRoutes[index].regex == pathPieces[1] {
			return node.subRoutes[index].search(pathPieces[1:])
		}
	}

	return NewRoute(http.MethodGet, "", http.NotFoundHandler())
}
