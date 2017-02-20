package router

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/xianggong/m2svis/server/modules/database"
)

// index to serve home page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../client/dist/index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func filterFormToSQL(r *http.Request) string {
	// Get filter information
	r.ParseForm()
	st := r.Form.Get("st")
	fn := r.Form.Get("fn")
	cu := r.Form.Get("cu")
	wf := r.Form.Get("wf")
	wg := r.Form.Get("wg")
	page := r.Form.Get("page")
	pagesize := r.Form.Get("pagesize")

	var filter []string
	if st != "" {
		filter = append(filter, "st>="+st)
	}
	if fn != "" {
		filter = append(filter, "fn<="+fn)
	}
	if cu != "" {
		filter = append(filter, "cu="+cu)
	}
	if wf != "" {
		filter = append(filter, "wf="+wf)
	}
	if wg != "" {
		filter = append(filter, "wg="+wg)
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

	// Pagination
	pageVal, _ := strconv.Atoi(page)
	pagesizeVal, _ := strconv.Atoi(pagesize)
	startIndex := (pageVal - 1) * pagesizeVal
	filterQuery += " LIMIT " + strconv.Itoa(startIndex) + "," + pagesize

	return filterQuery
}

func traceJSONHandler(w http.ResponseWriter, r *http.Request) {
	// Name and filter
	vars := mux.Vars(r)
	traceName := vars["traceName"]
	filterQuery := filterFormToSQL(r)

	// Query database
	data, err := database.GetInstTable(traceName, filterQuery)
	if err != nil {
		glog.Error(err)
		return
	}

	// Encode to JSON and write
	enc, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

func allHandler(w http.ResponseWriter, r *http.Request) {
	traces, err := database.GetTraces()
	if err != nil {
		glog.Error(err)
		return
	}

	enc, _ := json.Marshal(traces)
	w.Header().Set("Content-Type", "application/json")
	w.Write(enc)
}

// Init register handlers
func Init() {
	// Create router to redirect requests to handlers
	router := mux.NewRouter().StrictSlash(true)

	// Routes consist of a path and a handler function.
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/api/v1/traces/json", allHandler)
	router.HandleFunc("/api/v1/traces/{traceName}/json", traceJSONHandler)

	// Routes css/js/imgage
	fileServer := http.FileServer(http.Dir("../client/dist/"))
	router.PathPrefix("/static/").Handler(fileServer)

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":8000", router))
}
