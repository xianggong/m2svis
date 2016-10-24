package main

import (
	"fmt"
	"net/http"
)

func overview(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Fprintln(w, r.Form)
}

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/overview", overview)
	server.ListenAndServe()
}
