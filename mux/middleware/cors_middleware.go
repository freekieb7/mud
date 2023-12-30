package middleware

import (
	"net/http"
	"strconv"
	"strings"
)

type CorsMiddleware struct {
	allowedOrigins   []string
	allowedMethods   []string
	allowedHeaders   []string
	allowCredentials bool
	maxAge           int
}

func NewCorsMiddleware() Middleware {
	return &CorsMiddleware{
		allowedOrigins:   []string{"*"},
		allowedMethods:   []string{"*"},
		allowedHeaders:   []string{"*"},
		allowCredentials: false,
		maxAge:           86400,
	}
}

func (middleware CorsMiddleware) Process(next http.Handler) http.Handler {
	fn := func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Access-Control-Allow-Origin", strings.Join(middleware.allowedOrigins, ","))
		response.Header().Set("Access-Control-Allow-Methods", strings.Join(middleware.allowedMethods, ","))
		response.Header().Set("Access-Control-Allow-Headers", strings.Join(middleware.allowedHeaders, ","))
		response.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(middleware.allowCredentials))
		response.Header().Set("Access-Control-Max-Age", string(rune(middleware.maxAge)))

		next.ServeHTTP(response, request)
	}

	return http.HandlerFunc(fn)
}
