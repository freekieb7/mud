package mux

import (
	"log"
	"net/http"
	"strings"
)

type RouteTree struct {
	root RouteNode
}

type RouteNode struct {
	regex     string
	routes    []Route
	subRoutes map[string]*RouteNode
}

func (tree *RouteTree) Insert(route Route) {
	pathPieces := strings.Split(route.Path(), "/")
	tree.root.insert(route, pathPieces)
}

func (node *RouteNode) insert(route Route, pathPieces []string) {
	// Bad
	if node.regex != pathPieces[0] {
		log.Fatal("tree insertion: no path pieces left to assert")
	}

	// Match made in heaven
	if len(pathPieces) == 1 {
		node.routes = append(node.routes, route)
		return
	}

	// Maybe subroute matches
	if subRoute, ok := node.subRoutes[pathPieces[1]]; ok {
		subRoute.insert(route, pathPieces[1:])
		return
	}

	// No subroute matches, so create new matching subroute
	node.subRoutes[pathPieces[1]] = &RouteNode{
		regex:     pathPieces[1],
		routes:    []Route{},
		subRoutes: map[string]*RouteNode{},
	}

	// Continue insert with subroute
	if subRoute, ok := node.subRoutes[pathPieces[1]]; ok {
		subRoute.insert(route, pathPieces[1:])
		return
	}

	// Whut
	log.Fatal("tree insertion: failed subroute assertion after creation")
}

func (tree *RouteTree) Search(request *http.Request) Route {
	pathPieces := strings.Split(request.URL.Path, "/")
	return tree.root.search(pathPieces)
}

func (node *RouteNode) search(pathPieces []string) Route {
	// Match made in heaven
	if len(pathPieces) == 1 {
		if node.regex != pathPieces[0] {
			return NewNotFoundRoute()
		}

		return node.routes[0]
	}

	// We go the subroute route
	if subRoute, ok := node.subRoutes[pathPieces[1]]; ok {
		return subRoute.search(pathPieces[1:])
	}

	return NewRoute(http.MethodGet, "", http.NotFoundHandler())
}
