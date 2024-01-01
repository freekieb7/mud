package mux

import (
	"log"
	"net/http"
	"regexp"
	"strings"
)

type RouteTree interface {
	Insert(route Route)
	Search(request *http.Request) Route
	Routes() []Route
}

type routeTree struct {
	root routeNode
}

func NewRouteTree() RouteTree {
	return &routeTree{
		root: routeNode{
			regex:     "",
			routes:    []Route{},
			subRoutes: map[string]*routeNode{},
		},
	}
}

type routeNode struct {
	regex     string
	routes    []Route
	subRoutes map[string]*routeNode
}

func (tree *routeTree) Insert(route Route) {
	pathPieces := strings.Split(route.Path(), "/")
	tree.root.insert(route, pathPieces)
}

func (tree *routeTree) Search(request *http.Request) Route {
	pathPieces := strings.Split(request.URL.Path, "/")
	return tree.root.search(pathPieces, request.Method)
}

func (tree *routeTree) Routes() []Route {
	return tree.root.collectRoutes()
}

func (node *routeNode) collectRoutes() []Route {
	routes := node.routes

	for _, subRoutes := range node.subRoutes {
		routes = append(routes, subRoutes.collectRoutes()...)
	}

	return routes
}

func (node *routeNode) insert(route Route, pathPieces []string) {
	// Bad
	if node.regex != pathPieces[0] {
		log.Fatal("tree insertion: no path pieces left to assert")
	}

	// Match made in heaven
	if len(pathPieces) == 1 {
		node.routes = append(node.routes, route)
		return
	}

	// Maybe sub-route matches
	if subRoute, ok := node.subRoutes[pathPieces[1]]; ok {
		subRoute.insert(route, pathPieces[1:])
		return
	}

	// No sub-route matches, so create new matching sub-route
	node.subRoutes[pathPieces[1]] = &routeNode{
		regex:     pathPieces[1],
		routes:    []Route{},
		subRoutes: map[string]*routeNode{},
	}

	// Continue insert with sub-route
	if subRoute, ok := node.subRoutes[pathPieces[1]]; ok {
		subRoute.insert(route, pathPieces[1:])
		return
	}

	// WTF
	log.Fatal("tree insertion: failed subroute assertion after creation")
}

func (node *routeNode) search(pathPieces []string, method string) Route {
	// Match made in heaven
	if len(pathPieces) == 1 {
		// Check routes and methods
		for _, nodeRoute := range node.routes {
			if method == nodeRoute.Method() {
				return nodeRoute
			}
		}

		// No routes available
		if len(node.routes) == 0 {
			return NewNotFoundRoute()
		}

		return NewMethodNotAllowedRoute()
	}

	// We go the sub-route route
	for regex, subRoute := range node.subRoutes {
		// Static match
		if regex == pathPieces[1] {
			return subRoute.search(pathPieces[1:], method)
		}

		// Dynamic match
		if DynamicArgumentRegex.MatchString(regex) {
			regexParts := strings.Split(regex[1:len(regex)-1], ":")

			// Regex check
			if len(regexParts) == 2 {
				if regexp.MustCompile(regexParts[1]).MatchString(pathPieces[1]) {
					return subRoute.search(pathPieces[1:], method)
				}

				continue
			}

			// Only dynamic
			return subRoute.search(pathPieces[1:], method)
		}
	}

	return NewNotFoundRoute()
}
