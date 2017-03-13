package database

import "time"

// TraceInfo contains all traces information in database
type TraceInfo struct {
	Name  string    `db:"TraceName" json:"table_name"`
	Size  int       `db:"InstCount" json:"table_rows"`
	Stamp time.Time `db:"CreateTime" json:"create_time"`
}
