package trace

import (
	"errors"
	"log"
	"strconv"
)

// Activity of an instruction
type Activity struct {
	cycle    int
	activity string
}

// Instruction contains statistics of an instruction
type Instruction struct {
	Start              int    `db:"st"`
	Finish             int    `db:"fn"`
	Length             int    `db:"len"`
	FetchStart         int    `db:"fs"`
	FetchEnd           int    `db:"fe"`
	FetchStallWidth    int    `db:"fsw"`
	FetchStallBuffer   int    `db:"fsb"`
	IssueStart         int    `db:"is"`
	IssueEnd           int    `db:"ie"`
	IssieStallMax      int    `db:"ism"`
	IssieStallWidth    int    `db:"isw"`
	IssieStallBuffer   int    `db:"isb"`
	ReadStart          int    `db:"rs"`
	ReadEnd            int    `db:"re"`
	DecodeStart        int    `db:"ds"`
	DecodeEnd          int    `db:"de"`
	DecodeStallWidth   int    `db:"dsw"`
	DecodeStallBuffer  int    `db:"dsb"`
	ExecuteStart       int    `db:"es"`
	ExecuteEnd         int    `db:"ee"`
	ExecuteStallWidth  int    `db:"esw"`
	ExecuteStallBuffer int    `db:"esb"`
	WriteStart         int    `db:"ws"`
	WriteEnd           int    `db:"we"`
	WriteStallWidth    int    `db:"wsw"`
	WriteStallBuffer   int    `db:"wsb"`
	ID                 int    `db:"id"`
	CU                 int    `db:"cu"`
	IB                 int    `db:"ib"`
	WF                 int    `db:"wf"`
	WG                 int    `db:"wg"`
	UOP                int    `db:"uop"`
	ExecutionUnit      string `db:"eu"`
	Assembly           string `db:"asm"`
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
	for _, activity := range inst.LifeConcise {
		switch activity.activity {
		case "f":
			inst.FetchEnd = inst.FetchStart + activity.cycle
		case "s_cu_fe_rdy":
			inst.FetchStallBuffer = activity.cycle
			inst.FetchEnd += activity.cycle
		case "i":
			inst.IssueStart = inst.FetchEnd
		}
	}

	// Remove the last "end" activity
	inst.LifeConcise = inst.LifeConcise[:len(inst.LifeConcise)-1]
}

// InstructionJSON for representing instruction to timeline
type InstructionJSON struct {
	ID            string `json:"id"`
	Group         int    `json:"group"`
	Content       string `json:"content"`
	Start         int    `json:"start"`
	End           int    `json:"end"`
	SubGroup      string `json:"subgroup"`
	SubGroupOrder int    `json:"subgroupOrder"`
}

// GetJSON returns JSON representation of the instruction
func (inst *Instruction) GetJSON() []*InstructionJSON {
	// Store return
	var instructionJSONArray []*InstructionJSON

	// Common field
	instID := "inst_" + strconv.Itoa(inst.ID) + "_"
	subgroup := instID

	cycle := inst.Start
	for index, activity := range inst.LifeConcise {
		id := instID + activity.activity
		group := inst.CU
		content := activity.activity
		activityStart := cycle
		activityEnd := cycle + activity.cycle

		instJSON := InstructionJSON{ID: id, Group: group, Content: content,
			Start: activityStart, End: activityEnd, SubGroup: subgroup,
			SubGroupOrder: index}
		instructionJSONArray = append(instructionJSONArray, &instJSON)

		cycle = activityEnd
	}

	return instructionJSONArray
}

// GetOverviewJSON returns JSON representation of the overview of a instruction
func (inst *Instruction) GetOverviewJSON() []*InstructionJSON {
	// Store return
	var instructionJSONArray []*InstructionJSON

	id := "inst_" + strconv.Itoa(inst.ID) + "_" + strconv.Itoa(inst.CU)
	group := inst.CU
	subgroup := id
	content := inst.Assembly

	instJSON := InstructionJSON{ID: id, Group: group, Content: content,
		Start: inst.Start, End: inst.Start + inst.Length, SubGroup: subgroup}
	instructionJSONArray = append(instructionJSONArray, &instJSON)

	return instructionJSONArray

}
