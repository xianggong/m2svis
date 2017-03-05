package database

import "time"

// TraceAll contains all traces information in database
type TraceAll struct {
	Name  string    `db:"table_name" json:"table_name"`
	Size  int       `db:"table_rows" json:"table_rows"`
	Stamp time.Time `db:"create_time" json:"create_time"`
}

// TraceCount contains count information of a trace
type TraceCount struct {
	Count int `db:"count(*)" json:"count"`
}

// TraceData contains all information of an instruction in database
type TraceData struct {
	Start              int    `db:"st"   json:"Start"`
	Finish             int    `db:"fn"   json:"Finish"`
	Length             int    `db:"len"  json:"Len"`
	FetchStart         int    `db:"fs"   json:"Fs"`
	FetchEnd           int    `db:"fe"   json:"Fe"`
	FetchStallWidth    int    `db:"fsw"  json:"Fsw"`
	FetchStallBuffer   int    `db:"fsb"  json:"Fsb"`
	IssueStart         int    `db:"isu"  json:"Is"`
	IssueEnd           int    `db:"ie"   json:"Ie"`
	IssueStallMax      int    `db:"ism"  json:"Ism"`
	IssueStallWidth    int    `db:"isw"  json:"Isw"`
	IssueStallBuffer   int    `db:"isb"  json:"Isb"`
	DecodeStart        int    `db:"ds"   json:"Ds"`
	DecodeEnd          int    `db:"de"   json:"De"`
	DecodeStallWidth   int    `db:"dsw"  json:"Dsw"`
	DecodeStallBuffer  int    `db:"dsb"  json:"Dsb"`
	ReadStart          int    `db:"rs"   json:"Rs"`
	ReadEnd            int    `db:"re"   json:"Re"`
	ReadStallWidth     int    `db:"rsw"  json:"Rsw"`
	ReadStallBuffer    int    `db:"rsb"  json:"Rsb"`
	ExecuteStart       int    `db:"es"   json:"Es"`
	ExecuteEnd         int    `db:"ee"   json:"Ee"`
	ExecuteStallWidth  int    `db:"esw"  json:"Esw"`
	ExecuteStallBuffer int    `db:"esb"  json:"Esb"`
	WriteStart         int    `db:"ws"   json:"Ws"`
	WriteEnd           int    `db:"we"   json:"We"`
	WriteStallWidth    int    `db:"wsw"  json:"Wsw"`
	WriteStallBuffer   int    `db:"wsb"  json:"Wsb"`
	GID                int    `db:"gid"  json:"GUID"`
	ID                 int    `db:"id"   json:"ID"`
	CU                 int    `db:"cu"   json:"CU"`
	IB                 int    `db:"ib"   json:"IB"`
	WF                 int    `db:"wf"   json:"WF"`
	WG                 int    `db:"wg"   json:"WG"`
	UOP                int    `db:"uop"  json:"UOP"`
	ExecutionUnit      int    `db:"eu"   json:"EU"`
	Assembly           string `db:"asm"  json:"Asm"`
}

// InstTimelineJSON contains instruction timeline information
type InstTimelineJSON struct {
}
