package mux

import (
	"net/http"
	"strings"
)

type pathMatcher struct{}

func NewPathMatcher() Matcher {
	return &pathMatcher{}
}

func (matcher *pathMatcher) Match(route Route, request *http.Request) bool {
	// STEP 1: Remove query part from param
	uriSlice := strings.Split(request.URL.Path, "/")
	pathSlice := strings.Split(route.Path(), "/")

	if len(pathSlice) != len(uriSlice) {
		return false
	}

	for position, pathArgument := range pathSlice {
		uriArgument := uriSlice[position]

		// Check for matching dynamic argument
		if DynamicArgumentRegex.MatchString(pathArgument) {
			// TODO check regex of argument
			continue
		}

		// Check for matching static argument
		if pathArgument != uriArgument {
			return false
		}
	}

	return true
}
