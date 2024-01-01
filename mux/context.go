package mux

import (
	"net/http"
	"strings"
)

var (
	RouteCtxKey = "RouteContext"
)

func PathParam(request *http.Request, name string) string {
	currentRoute := request.Context().Value(RouteCtxKey).(Route)

	requestSlice := strings.Split(request.URL.Path, "/")
	pathSlice := strings.Split(currentRoute.Path(), "/")

	param := "{" + name + "}"

	for position, pathEntry := range pathSlice {
		if pathEntry == param {
			return requestSlice[position]
		}
	}

	return ""
}

func QueryParam(request *http.Request, name string) string {
	return request.URL.Query().Get(name)
}
