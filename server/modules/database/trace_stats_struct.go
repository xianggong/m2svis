package database

// TraceCount contains count information of a trace
type TraceCount struct {
	Count int `db:"CountInsts" json:"count"`
}

// TraceMeta contains overview information of a trace
type TraceMeta struct {
	CountInsts int `db:"CountInsts" json:"CountInsts"`
	MinCycle   int `db:"MinCycle" json:"MinCycle"`
	MaxCycle   int `db:"MaxCycle" json:"MaxCycle"`
	CountStall int `db:"CountStall" json:"CountStall"`
	CountWF    int `db:"CountWF" json:"CountWF"`
	CountWG    int `db:"CountWG" json:"CountWG"`
	CountCU    int `db:"CountCU" json:"CountCU"`
}

// TraceStall contains stall information
type TraceStall struct {
	StallTotal    int `db:"StallTotal" json:"StallTotal"`
	StallFrontend int `db:"StallFrontend" json:"StallFrontend"`
	StallIssue    int `db:"StallIssue" json:"StallIssue"`
	StallDecode   int `db:"StallDecode" json:"StallDecode"`
	StallRead     int `db:"StallRead" json:"StallRead"`
	StallExecute  int `db:"StallExecute" json:"StallExecute"`
	StallWrite    int `db:"StallWrite" json:"StallWrite"`
}

// TraceStallColumn contains stall information in columns
type TraceStallColumn struct {
	Frontend []int `db:"StallFrontend" json:"StallFrontend"`
	Issue    []int `db:"StallIssue" json:"StallIssue"`
	Decode   []int `db:"StallDecode" json:"StallDecode"`
	Read     []int `db:"StallRead" json:"StallRead"`
	Execute  []int `db:"StallExecute" json:"StallExecute"`
	Write    []int `db:"StallWrite" json:"StallWrite"`
}
