package middleware

import "net/http"

type Middleware interface {
	Process(next http.Handler) http.Handler
}
