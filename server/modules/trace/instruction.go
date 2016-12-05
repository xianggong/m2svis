package trace

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
)

// Activity of an instruction
type Activity struct {
	cycle    int
	activity string
}

// Instruction contains statistics of an instruction
type Instruction struct {
	Start              int    `db:"st" sql:"st INTEGER," json:"start"`
	Finish             int    `db:"fn" sql:"fn INTEGER," json:"end"`
	Length             int    `db:"len" sql:"len INTEGER,"`
	FetchStart         int    `db:"fs" sql:"fs INTEGER,"`
	FetchEnd           int    `db:"fe" sql:"fe INTEGER,"`
	FetchStallWidth    int    `db:"fsw" sql:"fsw INTEGER,"`
	FetchStallBuffer   int    `db:"fsb" sql:"fsb INTEGER,"`
	IssueStart         int    `db:"is" sql:"is INTEGER,"`
	IssueEnd           int    `db:"ie" sql:"ie INTEGER,"`
	IssueStallMax      int    `db:"ism" sql:"ism INTEGER,"`
	IssueStallWidth    int    `db:"isw" sql:"isw INTEGER,"`
	IssueStallBuffer   int    `db:"isb" sql:"isb INTEGER,"`
	ReadStart          int    `db:"rs" sql:"rs INTEGER,"`
	ReadEnd            int    `db:"re" sql:"re INTEGER,"`
	ReadStallWidth     int    `db:"rsw" sql:"rsw INTEGER,"`
	REadStallBuffer    int    `db:"rsb" sql:"rsb INTEGER,"`
	DecodeStart        int    `db:"ds" sql:"ds INTEGER,"`
	DecodeEnd          int    `db:"de" sql:"de INTEGER,"`
	DecodeStallWidth   int    `db:"dsw" sql:"dsw INTEGER,"`
	DecodeStallBuffer  int    `db:"dsb" sql:"dsb INTEGER,"`
	ExecuteStart       int    `db:"es" sql:"es INTEGER,"`
	ExecuteEnd         int    `db:"ee" sql:"ee INTEGER,"`
	ExecuteStallWidth  int    `db:"esw" sql:"esw INTEGER,"`
	ExecuteStallBuffer int    `db:"esb" sql:"esb INTEGER,"`
	WriteStart         int    `db:"ws" sql:"ws INTEGER,"`
	WriteEnd           int    `db:"we" sql:"we INTEGER,"`
	WriteStallWidth    int    `db:"wsw" sql:"wsw INTEGER,"`
	WriteStallBuffer   int    `db:"wsb" sql:"wsb INTEGER,"`
	ID                 int    `db:"id" sql:"id INTEGER," json:"subgroup"`
	CU                 int    `db:"cu" sql:"cu INTEGER," json:"group"`
	IB                 int    `db:"ib" sql:"ib INTEGER,"`
	WF                 int    `db:"wf" sql:"wf INTEGER,"`
	WG                 int    `db:"wg" sql:"wg INTEGER,"`
	UOP                int    `db:"uop" sql:"uop INTEGER,"`
	ExecutionUnit      int    `db:"eu" sql:"eu INTEGER,"`
	Assembly           string `db:"asm" sql:"asm VARCHAR" json:"content"`
	LifeConcise        []Activity
	LifeVerbose        []Activity
}

func (inst *Instruction) sanityCheck(info map[string]string) error {
	// Sanity check
	id, _ := strconv.Atoi(info["id"])
	cu, _ := strconv.Atoi(info["cu"])
	if id != inst.ID && cu != inst.CU {
		log.Printf("Expected id/cu=%d/%d, Actual id/cu=%d/%d\n",
			inst.ID, inst.CU, id, cu)
		return errors.New("Instruction: id/cu doesn't match!")
	}
	return nil
}

// New record 'New' activity of an instruction
func (inst *Instruction) New(cycle int, info map[string]string) {
	inst.Start = cycle
	inst.FetchStart = cycle
	inst.ID, _ = strconv.Atoi(info["id"])
	inst.CU, _ = strconv.Atoi(info["cu"])
	inst.IB, _ = strconv.Atoi(info["ib"])
	inst.WF, _ = strconv.Atoi(info["wf"])
	inst.WG, _ = strconv.Atoi(info["wg"])
	inst.UOP, _ = strconv.Atoi(info["uop_id"])
	inst.Assembly = info["asm"]
	inst.LifeVerbose = append(inst.LifeVerbose, Activity{cycle, info["stg"]})
	inst.LifeConcise = append(inst.LifeConcise, Activity{0, info["stg"]})
}

// Exe record 'Execute' activity of an instruction
func (inst *Instruction) Exe(cycle int, info map[string]string) {
	// Sanity check
	err := inst.sanityCheck(info)
	if err != nil {
		return
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
}

// End record 'End' activity of an instruction
func (inst *Instruction) End(cycle int, info map[string]string) {
	// Sanity check
	err := inst.sanityCheck(info)
	if err != nil {
		return
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
		case "bu-e", "lds-e", "su-e", "simd-e", "mem-e":
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
		// Stall: execute width
		case "s_br_wr_wth", "s_lds_wr_wth", "s_sl_wr_wth", "s_simd_wr_wth", "s_vc_mem_wr_wth":
			inst.WriteStallWidth += activity.cycle
			inst.WriteEnd = currEnd
		// Stall: execute buffer full
		case "s_br_wr_buf", "s_lds_wr_buf", "s_sl_wr_buf", "s_simd_wr_buf", "s_vc_mem_wr_buf":
			inst.WriteStallBuffer += activity.cycle
			inst.WriteEnd = currEnd
		}
	}

	// Remove the last "end" activity
	inst.LifeConcise = inst.LifeConcise[:len(inst.LifeConcise)-1]
}

// GetSQLQueryInsertTable returns SQL query string to insert an instruction table
func (inst *Instruction) GetSQLQueryInsertTable(tableName string) string {
	query := "CREATE TABLE " + tableName + "("

	// Loop through the struct's tags and append to query
	val := reflect.ValueOf(inst).Elem()
	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		if tag != "" {
			query += tag.Get("sql")
		}
	}
	query += ");"

	return query
}

// Dump returns a formatted string that is easy to read
func (inst *Instruction) Dump() string {
	var infoField string
	var infoValue string

	// Loop through the struct's tags and append to info
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

// // InstructionJSON for representing instruction to timeline
// type InstructionJSON struct {
// 	ID            string `json:"id"`
// 	Group         int    `json:"group"`
// 	Content       string `json:"content"`
// 	Start         int    `json:"start"`
// 	End           int    `json:"end"`
// 	SubGroup      string `json:"subgroup"`
// 	SubGroupOrder int    `json:"subgroupOrder"`
// }

// // GetJSON returns JSON representation of the instruction
// func (inst *Instruction) GetJSON() []*InstructionJSON {
// 	// Store return
// 	var instructionJSONArray []*InstructionJSON

// 	// Common field
// 	instID := "inst_" + strconv.Itoa(inst.ID) + "_"
// 	subgroup := instID

// 	cycle := inst.Start
// 	for index, activity := range inst.LifeConcise {
// 		id := instID + activity.activity
// 		group := inst.CU
// 		content := activity.activity
// 		activityStart := cycle
// 		activityEnd := cycle + activity.cycle

// 		instJSON := InstructionJSON{ID: id, Group: group, Content: content,
// 			Start: activityStart, End: activityEnd, SubGroup: subgroup,
// 			SubGroupOrder: index}
// 		instructionJSONArray = append(instructionJSONArray, &instJSON)

// 		cycle = activityEnd
// 	}

// 	return instructionJSONArray
// }

// // GetOverviewJSON returns JSON representation of the overview of a instruction
// func (inst *Instruction) GetOverviewJSON() []*InstructionJSON {
// 	// Store return
// 	var instructionJSONArray []*InstructionJSON

// 	id := "inst_" + strconv.Itoa(inst.ID) + "_" + strconv.Itoa(inst.CU)
// 	group := inst.CU
// 	subgroup := id
// 	content := inst.Assembly

// 	instJSON := InstructionJSON{ID: id, Group: group, Content: content,
// 		Start: inst.Start, End: inst.Start + inst.Length, SubGroup: subgroup}
// 	instructionJSONArray = append(instructionJSONArray, &instJSON)

// 	return instructionJSONArray

// }
