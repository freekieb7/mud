package main

import (
	"github.com/freekieb7/mud/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	router.Get("/test/test/asd", PathCheckHandler)

	log.Fatal(http.ListenAndServe(":3000", router))
}

func PathCheckHandler(response http.ResponseWriter, request *http.Request) {
	response.Write([]byte("Path: " + request.URL.Path))
}
