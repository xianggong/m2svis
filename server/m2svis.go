package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/xianggong/m2svis/server/modules/trace"
)

// Index To serve home page
func Index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintln(w, r.Form)
}

// Timeline To serve /timeline page
func Timeline(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/timeline.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// TimelineJSON To return timeline json data
func TimelineJSON(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// start := r.Form["start"]
	// end := r.Form["end"]
	w.Header().Set("Content-Type", "application/json")

	trace := new(trace.Trace)
	trace.Process("modules/trace/test.gz")

	enc, _ := json.Marshal(trace.GetJSON())

	w.Write(enc)
}

func main() {
	// Create router to redirect requests to handlers
	router := mux.NewRouter()

	// Routes consist of a path and a handler function.
	router.HandleFunc("/", Index)
	router.HandleFunc("/timeline", Timeline)
	router.HandleFunc("/timeline/json", TimelineJSON)

	// Routs css/js/imgage etc to ../client directory
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../client/")))

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", router))
}
