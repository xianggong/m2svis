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
	router.HandleFunc("/api/v1/traces/{traceName}/meta", traceMeta)
	router.HandleFunc("/api/v1/traces/{traceName}/stall", traceStall)
	router.HandleFunc("/api/v1/traces/{traceName}/stall/row", traceStallRow)
	router.HandleFunc("/api/v1/traces/{traceName}/stall/column", traceStallColumn)
	router.HandleFunc("/api/v1/traces/{traceName}/{cuid}/instruction/active", CUActiveInsts)

	// Routes css/js/images
	fileServer := http.FileServer(http.Dir("../client/dist/"))
	router.PathPrefix("/static/").Handler(fileServer)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", router))
}
