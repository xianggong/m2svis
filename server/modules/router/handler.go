package router

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/xianggong/m2svis/server/modules/database"
)

// index to serve home page
func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../client/dist/index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func traceData(w http.ResponseWriter, r *http.Request) {
	// Name and filter
	vars := mux.Vars(r)
	traceName := vars["traceName"]
	filterQuery := formToSQL(r)

	// Query database
	data, err := database.GetTraceData(traceName, filterQuery)
	if err != nil {
		glog.Error(err)
		return
	}

	// Encode to JSON and write
	enc, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

func traceAll(w http.ResponseWriter, r *http.Request) {
	traces, err := database.GetTraceAll()
	if err != nil {
		glog.Error(err)
		return
	}

	enc, _ := json.Marshal(traces)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

func traceCount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceName := vars["traceName"]
	filterQuery := formFiltersToSQL(r)

	traceCount, err := database.GetTraceCount(traceName, filterQuery)
	if err != nil {
		glog.Error(err)
		return
	}

	// Encode to JSON and write
	enc, _ := json.Marshal(traceCount)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}
