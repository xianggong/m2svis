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

type InstructionCSV struct {
	Start         int    `db:"Start"         dtype:"INTEGER" json:"Start"`
	Finish        int    `db:"Finish"        dtype:"INTEGER" json:"Finish"`
	Length        int    `db:"Length"        dtype:"INTEGER" json:"Length"`
	FetchBegin    int    `db:"FetchBegin"    dtype:"INTEGER" json:"FetchBegin"`
	FetchStall    int    `db:"FetchStall"    dtype:"INTEGER" json:"FetchStall"`
	FetchActive   int    `db:"FetchActive"   dtype:"INTEGER" json:"FetchActive"`
	FetchEnd      int    `db:"FetchEnd"      dtype:"INTEGER" json:"FetchEnd"`
	IssueBegin    int    `db:"IssueBegin"    dtype:"INTEGER" json:"IssueBegin"`
	IssueStall    int    `db:"IssueStall"    dtype:"INTEGER" json:"IssueStall"`
	IssueActive   int    `db:"IssueActive"   dtype:"INTEGER" json:"IssueActive"`
	IssueEnd      int    `db:"IssueEnd"      dtype:"INTEGER" json:"IssueEnd"`
	DecodeBegin   int    `db:"DecodeBegin"   dtype:"INTEGER" json:"DecodeBegin"`
	DecodeStall   int    `db:"DecodeStall"   dtype:"INTEGER" json:"DecodeStall"`
	DecodeActive  int    `db:"DecodeActive"  dtype:"INTEGER" json:"DecodeActive"`
	DecodeEnd     int    `db:"DecodeEnd"     dtype:"INTEGER" json:"DecodeEnd"`
	ReadBegin     int    `db:"ReadBegin"     dtype:"INTEGER" json:"ReadBegin"`
	ReadStall     int    `db:"ReadStall"     dtype:"INTEGER" json:"ReadStall"`
	ReadActive    int    `db:"ReadActive"    dtype:"INTEGER" json:"ReadActive"`
	ReadEnd       int    `db:"ReadEnd"       dtype:"INTEGER" json:"ReadEnd"`
	ExecuteBegin  int    `db:"ExecuteBegin"  dtype:"INTEGER" json:"ExecuteBegin"`
	ExecuteStall  int    `db:"ExecuteStall"  dtype:"INTEGER" json:"ExecuteStall"`
	ExecuteActive int    `db:"ExecuteActive" dtype:"INTEGER" json:"ExecuteActive"`
	ExecuteEnd    int    `db:"ExecuteEnd"    dtype:"INTEGER" json:"ExecuteEnd"`
	WriteBegin    int    `db:"WriteBegin"    dtype:"INTEGER" json:"WriteBegin"`
	WriteStall    int    `db:"WriteStall"    dtype:"INTEGER" json:"WriteStall"`
	WriteActive   int    `db:"WriteActive"   dtype:"INTEGER" json:"WriteActive"`
	WriteEnd      int    `db:"WriteEnd"      dtype:"INTEGER" json:"WriteEnd"`
	GID           int    `db:"GID"           dtype:"INTEGER" json:"GID"`
	ID            int    `db:"ID"            dtype:"INTEGER" json:"ID"`
	CU            int    `db:"CU"            dtype:"INTEGER" json:"CU"`
	IB            int    `db:"IB"            dtype:"INTEGER" json:"IB"`
	WF            int    `db:"WF"            dtype:"INTEGER" json:"WF"`
	WG            int    `db:"WG"            dtype:"INTEGER" json:"WG"`
	UOP           int    `db:"UOP"           dtype:"INTEGER" json:"UOP"`
	ExecutionUnit string `db:"ExecutionUnit" dtype:"TEXT"    json:"ExecutionUnit"`
	Type          string `db:"Type"          dtype:"TEXT"    json:"Type"`
	Assembly      string `db:"Assembly"      dtype:"TEXT"    json:"Assembly"`
}
