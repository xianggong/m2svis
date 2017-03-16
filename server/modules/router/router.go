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
	router.HandleFunc("/api/v1/traces/{traceName}/rawdata/data", rawdata)
	router.HandleFunc("/api/v1/traces/{traceName}/count", count)
	router.HandleFunc("/api/v1/traces/{traceName}/meta", meta)
	router.HandleFunc("/api/v1/traces/{traceName}/stats/stall", statsStall)
	router.HandleFunc("/api/v1/traces/{traceName}/stats/stall/row", statsStallRow)
	router.HandleFunc("/api/v1/traces/{traceName}/stats/stall/column", statsStallColumn)
	router.HandleFunc("/api/v1/traces/{traceName}/insts/active", CUActiveInsts)
	router.HandleFunc("/api/v1/traces/{traceName}/insts/count/insttype", instCountByInstType)
	router.HandleFunc("/api/v1/traces/{traceName}/insts/count/execunit", instCountByExecUnit)
	router.HandleFunc("/api/v1/traces/{traceName}/cycle/count/insttype", cycleCountByInstType)
	router.HandleFunc("/api/v1/traces/{traceName}/cycle/count/execunit", cycleCountByExecUnit)
	router.HandleFunc("/api/v1/traces/{traceName}/cycle/count/cu", cycleCountByCU)

	// Routes css/js/images
	fileServer := http.FileServer(http.Dir("../client/dist/"))
	router.PathPrefix("/static/").Handler(fileServer)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", router))
}
