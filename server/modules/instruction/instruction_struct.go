package instruction

// Activity of an instruction
type Activity struct {
	cycle    int
	activity string
}

// Instruction contains statistics of an instruction
type Instruction struct {
	Start              int    `db:"st"  dtype:"INTEGER"`
	Finish             int    `db:"fn"  dtype:"INTEGER"`
	Length             int    `db:"len" dtype:"INTEGER"`
	FetchStart         int    `db:"fs"  dtype:"INTEGER"`
	FetchEnd           int    `db:"fe"  dtype:"INTEGER"`
	FetchStallWidth    int    `db:"fsw" dtype:"INTEGER"`
	FetchStallBuffer   int    `db:"fsb" dtype:"INTEGER"`
	IssueStart         int    `db:"isu" dtype:"INTEGER"`
	IssueEnd           int    `db:"ie"  dtype:"INTEGER"`
	IssueStallMax      int    `db:"ism" dtype:"INTEGER"`
	IssueStallWidth    int    `db:"isw" dtype:"INTEGER"`
	IssueStallBuffer   int    `db:"isb" dtype:"INTEGER"`
	DecodeStart        int    `db:"ds"  dtype:"INTEGER"`
	DecodeEnd          int    `db:"de"  dtype:"INTEGER"`
	DecodeStallWidth   int    `db:"dsw" dtype:"INTEGER"`
	DecodeStallBuffer  int    `db:"dsb" dtype:"INTEGER"`
	ReadStart          int    `db:"rs"  dtype:"INTEGER"`
	ReadEnd            int    `db:"re"  dtype:"INTEGER"`
	ReadStallWidth     int    `db:"rsw" dtype:"INTEGER"`
	ReadStallBuffer    int    `db:"rsb" dtype:"INTEGER"`
	ExecuteStart       int    `db:"es"  dtype:"INTEGER"`
	ExecuteEnd         int    `db:"ee"  dtype:"INTEGER"`
	ExecuteStallWidth  int    `db:"esw" dtype:"INTEGER"`
	ExecuteStallBuffer int    `db:"esb" dtype:"INTEGER"`
	WriteStart         int    `db:"ws"  dtype:"INTEGER"`
	WriteEnd           int    `db:"we"  dtype:"INTEGER"`
	WriteStallWidth    int    `db:"wsw" dtype:"INTEGER"`
	WriteStallBuffer   int    `db:"wsb" dtype:"INTEGER"`
	GID                int    `db:"gid" dtype:"INTEGER"`
	ID                 int    `db:"id"  dtype:"INTEGER"`
	CU                 int    `db:"cu"  dtype:"INTEGER"`
	IB                 int    `db:"ib"  dtype:"INTEGER"`
	WF                 int    `db:"wf"  dtype:"INTEGER"`
	WG                 int    `db:"wg"  dtype:"INTEGER"`
	UOP                int    `db:"uop" dtype:"INTEGER"`
	Assembly           string `db:"asm" dtype:"TEXT"`
	ExecutionUnit      string `db:"eu"  dtype:"TEXT"`
	Type               string `db:"type" dtype:"TEXT"`
	LifeConcise        []Activity
	LifeVerbose        []Activity
}
