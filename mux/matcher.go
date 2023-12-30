package mux

import (
	"net/http"
)

type Matcher interface {
	Match(route Route, request *http.Request) bool
}
