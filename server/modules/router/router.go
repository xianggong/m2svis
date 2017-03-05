package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Init register handlers
func Init() {
	// Create router to redirect requests to handlers
	router := mux.NewRouter().StrictSlash(true)

	// Routes consist of a path and a handler function.
	router.HandleFunc("/", index)
	router.HandleFunc("/api/v1/traces/data", traceAll)
	router.HandleFunc("/api/v1/traces/{traceName}/data", traceData)
	router.HandleFunc("/api/v1/traces/{traceName}/count", traceCount)

	// Routes css/js/images
	fileServer := http.FileServer(http.Dir("../client/dist/"))
	router.PathPrefix("/static/").Handler(fileServer)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", router))
}
