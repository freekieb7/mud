package middleware

import (
	"log"
	"net/http"
)

type bullshitMiddleware struct {
}

func NewBullshitMiddleware() Middleware {
	return &bullshitMiddleware{}
}

func (middleware *bullshitMiddleware) Process(next http.Handler) http.Handler {
	fn := func(response http.ResponseWriter, request *http.Request) {
		next.ServeHTTP(response, request)
		log.Println("bulshit")
	}

	return http.HandlerFunc(fn)
}
