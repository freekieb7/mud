package middleware

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type loggingMiddleware struct {
	writers []io.Writer
}

type LoggingMiddlewareConfig struct {
	Writers []io.Writer
}

func NewLoggingMiddleware(config LoggingMiddlewareConfig) Middleware {
	return &loggingMiddleware{
		writers: config.Writers,
	}
}

func NewDefaultLoggingMiddleware() Middleware {
	file, err := os.OpenFile("/tmp/test.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Fatalf("middleware: error opening log file: %v", err)
	}

	return NewLoggingMiddleware(
		LoggingMiddlewareConfig{
			Writers: []io.Writer{
				file,
				os.Stdout,
			},
		},
	)
}

func (middleware *loggingMiddleware) Process(next http.Handler) http.Handler {
	fn := func(response http.ResponseWriter, request *http.Request) {
		now := time.Now()

		next.ServeHTTP(response, request)

		duration := time.Since(now)

		go middleware.logResult(duration)
	}

	return http.HandlerFunc(fn)
}

func (middleware *loggingMiddleware) logResult(duration time.Duration) {
	multi := io.MultiWriter(middleware.writers...)
	log.SetOutput(multi)

	log.Printf(
		"%s %s %s %s %s %s %s",
		time.Now(), // Date
		"-",        // Channel
		"-",        // Level
		"-",        // Message
		"-",        // Context
		"-",        // Extra
		duration,   // Duration
	)
}
