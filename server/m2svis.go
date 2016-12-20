package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/xianggong/m2svis/server/modules/trace"
)

// Index to serve home page
func Index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintln(w, r.Form)
}

// Timeline to serve /timeline page
func Timeline(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/timeline.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// TimelineJSON to return timeline json data
func TimelineJSON(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	traceName := vars["traceName"]

	// Get filter information
	r.ParseForm()
	st := r.Form["st"]
	ed := r.Form["fn"]
	cu := r.Form["cu"]
	wf := r.Form["wf"]
	wg := r.Form["wg"]

	var filter []string
	if len(st) != 0 {
		filter = append(filter, "st>="+st[0])
	}
	if len(ed) != 0 {
		filter = append(filter, "fn<="+ed[0])
	}
	if len(cu) != 0 {
		for _, id := range cu {
			filter = append(filter, "cu="+id)
		}
	}
	if len(wf) != 0 {
		for _, id := range wf {
			filter = append(filter, "wf="+id)
		}
	}
	if len(wg) != 0 {
		for _, id := range wf {
			filter = append(filter, "wg="+id)
		}
	}

	// Flatten to SQL query
	filterQuery := ""
	if len(filter) != 0 {
		filterQuery += " WHERE "
		for idx, val := range filter {
			filterQuery += val
			if idx != len(filter)-1 {
				filterQuery += " AND "
			}
		}
	}

	// Get data in JSON format
	data, err := trace.GetInstance().GetInstsInDB(traceName, filterQuery)
	if err != nil {
		glog.Error(err)
		return
	}
	enc, _ := json.Marshal(data)

	// Write JSON data
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

func initModules() {
	// Create trace module
	trace := trace.GetInstance()
	trace.Init("./m2svis.toml")
}

func initRouter() {
	// Create router to redirect requests to handlers
	router := mux.NewRouter().StrictSlash(true)

	// Routes consist of a path and a handler function.
	router.HandleFunc("/", Index)
	router.HandleFunc("/{traceName}/timeline", Timeline)
	router.HandleFunc("/{traceName}/timeline/json", TimelineJSON)

	// Routs css/js/imgage etc to ../client directory
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../client/")))

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", router))
}

func main() {

	// Init modules
	initModules()

	// Init routers
	initRouter()

}
