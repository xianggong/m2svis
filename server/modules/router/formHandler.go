package router

import (
	"fmt"
	"net/http"
	"strconv"
)

func formFiltersToSQL(r *http.Request) string {
	r.ParseForm()
	return r.Form.Get("filters")
}

func rawDataFormToSQL(r *http.Request) string {
	query := ""
	r.ParseForm()

	// Filters
	filters := r.Form.Get("filters")
	query += filters

	// Pagination
	page, _ := strconv.Atoi(r.Form.Get("page"))
	size, _ := strconv.Atoi(r.Form.Get("size"))
	if page != 0 && size != 0 {
		paginationQuery := fmt.Sprintf(" LIMIT %d, %d ", (page-1)*size, size)
		query += paginationQuery
	}

	return query
}
