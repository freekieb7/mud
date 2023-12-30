package mux

import (
	"net/http"
)

type methodMatcher struct{}

func NewMethodMatcher() Matcher {
	return &methodMatcher{}
}

func (matcher *methodMatcher) Match(route Route, request *http.Request) bool {
	return request.Method == route.Method()
}
