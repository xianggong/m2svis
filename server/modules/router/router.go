package router

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var routerMap = map[string]func(http.ResponseWriter, *http.Request){
	"/": index,
	"/api/v1/traces/data":                              traceAll,
	"/api/v1/traces/{traceName}/rawdata/data":          rawdata,
	"/api/v1/traces/{traceName}/count":                 count,
	"/api/v1/traces/{traceName}/meta":                  meta,
	"/api/v1/traces/{traceName}/stats/stall":           statsStall,
	"/api/v1/traces/{traceName}/stats/stall/column":    statsStallColumn,
	"/api/v1/traces/{traceName}/insts/active":          CUActiveInsts,
	"/api/v1/traces/{traceName}/insts/type":            instType,
	"/api/v1/traces/{traceName}/insts/count/insttype":  instCountByInstType,
	"/api/v1/traces/{traceName}/insts/count/execunit":  instCountByExecUnit,
	"/api/v1/traces/{traceName}/insts/length/insttype": instLengthByInstType,
	"/api/v1/traces/{traceName}/insts/length/execunit": instLengthByExecUnit,
	"/api/v1/traces/{traceName}/cycle/count/insttype":  cycleCountByInstType,
	"/api/v1/traces/{traceName}/cycle/count/execunit":  cycleCountByExecUnit,
	"/api/v1/traces/{traceName}/cycle/count/cu":        cycleCountByCU,
	"/api/v1/traces/{traceName}/execunit/utilization":  execUnitUtilization,
}

// Init register handlers
func Init() {
	// Create router to redirect requests to handlers
	router := mux.NewRouter().StrictSlash(true)

	// Register handlers
	for path, pathHandler := range routerMap {
		router.HandleFunc(path, pathHandler)
	}

	// Routes css/js/images
	fileServer := http.FileServer(http.Dir("../client/dist/"))
	router.PathPrefix("/static/").Handler(fileServer)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", router))
}
