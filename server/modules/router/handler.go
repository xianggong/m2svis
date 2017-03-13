package router

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

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
	data, err := database.GetAllTraceInfo()
	if err != nil {
		glog.Error(err)
		return
	}

	enc, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

func rawdata(w http.ResponseWriter, r *http.Request) {
	// Name and filter
	vars := mux.Vars(r)
	traceName := vars["traceName"]
	filterQuery := rawDataFormToSQL(r)

	// Query database
	data, err := database.GetTraceRawdata(traceName, filterQuery)
	if err != nil {
		glog.Error(err)
		return
	}

	// Encode to JSON and write
	enc, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

func count(w http.ResponseWriter, r *http.Request) {
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

func meta(w http.ResponseWriter, r *http.Request) {
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

func statsStall(w http.ResponseWriter, r *http.Request) {
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

func statsStallRow(w http.ResponseWriter, r *http.Request) {
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

func statsStallColumn(w http.ResponseWriter, r *http.Request) {
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

// CUActiveInsts returns active instructions
func CUActiveInsts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceName := vars["traceName"]

	r.ParseForm()

	start, _ := strconv.Atoi(r.Form.Get("start"))
	finish, _ := strconv.Atoi(r.Form.Get("finish"))
	cu, _ := strconv.Atoi(r.Form.Get("cu"))
	windowSize, _ := strconv.Atoi(r.Form.Get("windowSize"))

	activeInsts, err := database.GetTraceActiveCount(traceName, cu, start, finish, windowSize)
	if err != nil {
		glog.Error(err)
		return
	}

	// Encode to JSON and write
	enc, _ := json.Marshal(activeInsts)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

func instCountByInstType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceName := vars["traceName"]

	data := database.GetInstCountByInstType(traceName)

	// Encode to JSON and write
	enc, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

func instCountByExecUnit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceName := vars["traceName"]

	data := database.GetInstCountByExecUnit(traceName)

	// Encode to JSON and write
	enc, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}
