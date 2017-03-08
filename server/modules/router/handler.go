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

func traceMeta(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceName := vars["traceName"]

	traceMeta, err := database.GetTraceMeta(traceName, "")
	if err != nil {
		glog.Error(err)
		return
	}

	// Encode to JSON and write
	enc, _ := json.Marshal(traceMeta)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

func traceStall(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceName := vars["traceName"]

	traceStall, err := database.GetTraceStall(traceName, "")
	if err != nil {
		glog.Error(err)
		return
	}

	// Encode to JSON and write
	enc, _ := json.Marshal(traceStall)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

func traceStallRow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceName := vars["traceName"]

	stallRow, err := database.GetTraceStallRow(traceName, "")
	if err != nil {
		glog.Error(err)
		return
	}

	// Encode to JSON and write
	enc, _ := json.Marshal(stallRow)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

func traceStallColumn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceName := vars["traceName"]

	stallColumn, err := database.GetTraceStallColumn(traceName, "")
	if err != nil {
		glog.Error(err)
		return
	}

	// Encode to JSON and write
	enc, _ := json.Marshal(stallColumn)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

func CUActiveInsts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceName := vars["traceName"]

	activeInsts, err := database.GetTraceActiveCount(traceName, 0, 0, 10000, 0)
	if err != nil {
		glog.Error(err)
		return
	}

	// Encode to JSON and write
	enc, _ := json.Marshal(activeInsts)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}
