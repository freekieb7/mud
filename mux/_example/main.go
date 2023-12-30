package main

import (
	"fmt"
	"github.com/freekieb7/mud/mux"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"
)

func main() {
	router := mux.NewDefaultRouter()

	//router.Group("/", func(router mux.Router) {
	//	router.Get("/test", PathCheckHandler)
	//	router.Group("/b", func(router mux.Router) {
	//		router.Get("/test", PathCheckHandler)
	//	})
	//	router.Get("/testa", PathCheckHandler)
	//})
	//
	//router.Get("/", PathCheckHandler)
	router.Get("/test/{something}/asd", PathCheckHandler)

	logRoutes(router)

	log.Fatal(http.ListenAndServe(":3000", router))
}

func PathCheckHandler(response http.ResponseWriter, request *http.Request) {
	pParam := mux.PathParam(request, "something")
	qParam := mux.QueryParam(request, "test")

	//response.Write([]byte("Path: " + request.URL.Path))
	response.Write([]byte("Param: " + pParam + "\n"))
	response.Write([]byte("Query: " + qParam))
}

func logRoutes(router mux.Router) {
	time.Sleep(time.Millisecond * 100)

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", "Name", "Method", "Scheme", "Host", "Path")
	for _, route := range router.Routes() {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", "-", route.Method, "ANY", "ANY", route.Path)
	}
	w.Flush()
}
