package instruction

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

var instGID = 0

func getGID() int {
	instGID++
	return instGID
}

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
	ExecutionUnit      int    `db:"eu"  dtype:"INTEGER"`
	Assembly           string `db:"asm" dtype:"TEXT"`
	LifeConcise        []Activity
	LifeVerbose        []Activity
}

func (inst *Instruction) matchCheck(info map[string]string) error {
	// Sanity check
	id, _ := strconv.Atoi(info["id"])
	cu, _ := strconv.Atoi(info["cu"])
	if id != inst.ID && cu != inst.CU {
		log.Printf("Expected id/cu=%d/%d, Actual id/cu=%d/%d\n",
			inst.ID, inst.CU, id, cu)
		return errors.New("id/cu does not match")
	}

	// Return
	return nil
}

// New record 'New' activity of an instruction
func (inst *Instruction) New(cycle int, info map[string]string) error {
	// Record statistics
	inst.Start = cycle
	inst.FetchStart = cycle
	inst.GID = getGID()
	inst.ID, _ = strconv.Atoi(info["id"])
	inst.CU, _ = strconv.Atoi(info["cu"])
	inst.IB, _ = strconv.Atoi(info["ib"])
	inst.WF, _ = strconv.Atoi(info["wf"])
	inst.WG, _ = strconv.Atoi(info["wg"])
	inst.UOP, _ = strconv.Atoi(info["uop_id"])
	// Remove comment
	commentIndex := strings.Index(info["asm"], " //")
	inst.Assembly = "[" + info["wg"] + "-" + info["wf"] + "]: " + info["asm"][:commentIndex]
	inst.LifeVerbose = append(inst.LifeVerbose, Activity{cycle, info["stg"]})
	inst.LifeConcise = append(inst.LifeConcise, Activity{0, info["stg"]})

	// Return
	return nil
}

// Exe record 'Execute' activity of an instruction
func (inst *Instruction) Exe(cycle int, info map[string]string) error {
	// Sanity check
	err := inst.matchCheck(info)
	if err != nil {
		return err
	}

	// Get current/previous activity and time interval
	instPrevActivity := inst.LifeVerbose[len(inst.LifeVerbose)-1]
	instCurrActivity := Activity{cycle, info["stg"]}
	timeInterval := instCurrActivity.cycle - instPrevActivity.cycle

	// Record verbose activity
	inst.LifeVerbose = append(inst.LifeVerbose, instCurrActivity)

	// Update concise activity
	inst.LifeConcise[len(inst.LifeConcise)-1].cycle += timeInterval
	if instCurrActivity.activity != instPrevActivity.activity {
		activity := Activity{0, instCurrActivity.activity}
		inst.LifeConcise = append(inst.LifeConcise, activity)
	}

	// Return
	return nil
}

// End record 'End' activity of an instruction
func (inst *Instruction) End(cycle int, info map[string]string) error {
	// Sanity check
	err := inst.matchCheck(info)
	if err != nil {
		return err
	}

	// Update
	inst.Finish = cycle
	inst.Length = inst.Finish - inst.Start
	info["stg"] = "end"
	inst.Exe(cycle, info)

	// Update statistics
	currEnd := inst.Start
	for _, activity := range inst.LifeConcise {
		currStart := currEnd
		currEnd += activity.cycle
		switch activity.activity {
		// Fetch
		case "f":
			inst.FetchEnd = currEnd
		// Stall: fetch is not ready
		case "s_cu_fe_rdy":
			inst.FetchStallBuffer += activity.cycle
			inst.FetchEnd = currEnd
		// Issue
		case "i":
			inst.IssueStart = currStart
			inst.IssueEnd = currEnd
		// Stall: issue width
		case "s_cu_iss_wth":
			inst.IssueStallWidth += activity.cycle
			inst.IssueEnd = currEnd
		// Stall: max instruction already issued
		case "s_br_iss_max", "s_sl_iss_max", "s_lds_iss_max", "s_vc_iss_max":
			inst.IssueStallMax += activity.cycle
			inst.IssueEnd = currEnd
		// Stall: issue buffer full
		case "s_br_iss_buf", "s_sl_iss_buf", "s_lds_iss_buf", "s_vc_iss_buf":
			inst.IssueStallBuffer += activity.cycle
			inst.IssueEnd = currEnd
		// Read
		case "bu-r", "lds-r", "su-r", "simd-r", "mem-r":
			inst.ReadStart = currStart
			inst.ReadEnd = currEnd
		// Stall: read width
		case "s_br_rd_wth", "s_sl_rd_wth", "s_lds_rd_wth", "s_vc_mem_rd_wth":
			inst.ReadStallWidth += activity.cycle
			inst.ReadEnd = currEnd
		// Stall: read buffer full
		case "s_br_rd_buf", "s_sl_rd_buf", "s_lds_rd_buf", "s_vc_mem_rd_buf":
			inst.ReadStallWidth += activity.cycle
			inst.ReadEnd = currEnd
		// Decode
		case "bu-d", "lds-d", "su-d", "simd-d", "mem-d":
			inst.DecodeStart = currStart
			inst.DecodeEnd = currEnd
		// Stall: decode width
		case "s_br_dec_wth", "s_lds_dec_wth", "s_sl_dec_wth", "s_simd_dec_wth", "s_vc_mem_dec_wth":
			inst.DecodeStallWidth += activity.cycle
			inst.DecodeEnd = currEnd
		// Stall: decode buffer full
		case "s_br_dec_buf", "s_lds_dec_buf", "s_sl_dec_buf", "s_simd_dec_buf", "s_vc_mem_dec_buf":
			inst.DecodeStallBuffer += activity.cycle
			inst.DecodeEnd = currEnd
			// Execution
		case "bu-e", "lds-e", "su-e", "simd-e", "mem-e", "su-m", "mem-m":
			inst.ExecuteStart = currStart
			inst.ExecuteEnd = currEnd
		// Stall: execute width
		case "s_br_exe_wth", "s_lds_exe_wth", "s_sl_exe_wth", "s_simd_exe_wth", "s_vc_mem_exe_wth":
			inst.ExecuteStallWidth += activity.cycle
			inst.ExecuteEnd = currEnd
		// Stall: execute buffer full
		case "s_br_exe_buf", "s_lds_exe_buf", "s_sl_exe_buf", "s_simd_exe_buf", "s_vc_mem_exe_buf":
			inst.ExecuteStallBuffer += activity.cycle
			inst.ExecuteEnd = currEnd
		// Write
		case "bu-w", "lds-w", "su-w", "simd-w", "mem-w":
			inst.WriteStart = currStart
			inst.WriteEnd = currEnd
		// Stall: write width
		case "s_br_wr_wth", "s_lds_wr_wth", "s_sl_wr_wth", "s_simd_wr_wth", "s_vc_mem_wr_wth":
			inst.WriteStallWidth += activity.cycle
			inst.WriteEnd = currEnd
		// Stall: write buffer full
		case "s_br_wr_buf", "s_lds_wr_buf", "s_sl_wr_buf", "s_simd_wr_buf", "s_vc_mem_wr_buf":
			inst.WriteStallBuffer += activity.cycle
			inst.WriteEnd = currEnd
		}
	}

	// Remove the last "end" activity
	inst.LifeConcise = inst.LifeConcise[:len(inst.LifeConcise)-1]

	// Return
	return nil
}

// GetInstSQLColNames returns all columns name of the instruction struct
func GetInstSQLColNames(prefix string, suffix string) (str string) {
	str = "("

	// Loop through the struct's tags and append to query
	inst := new(Instruction)
	val := reflect.ValueOf(inst).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		if tag != "" {
			str += prefix + tag.Get("db") + suffix
		}
	}

	str = strings.TrimSuffix(str, suffix)
	str += ")"

	return str
}

// QueryCreateInstTable returns SQL query to insert an instruction table
func QueryCreateInstTable(tableName string) string {
	inst := new(Instruction)
	query := "CREATE TABLE IF NOT EXISTS " + tableName + "_insts ("

	// Loop through the struct's tags and append to query
	val := reflect.ValueOf(inst).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		if tag != "" {
			dbColName := tag.Get("db")
			dbColType := tag.Get("dtype")
			if dbColName != "-" && dbColType != "-" {
				query += dbColName + " " + dbColType + ", "
			}
		}
	}

	// Primary key
	query += " PRIMARY KEY (gid));"

	return query
}

// Dump returns a formatted string that is friendly to read
func (inst *Instruction) Dump() string {
	var infoField string
	var infoValue string

	// Use reflect to loop through the struct's tags and append to info
	val := reflect.ValueOf(inst).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		if tag != "" {
			infoField += fmt.Sprintf("%v\t", tag.Get("db"))
			infoValue += fmt.Sprintf("%v\t", valueField.Interface())
		}
	}
	return infoField + "\n" + infoValue
}
