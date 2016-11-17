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
	id                                   int
	start, finish, length                int
	fetch, decode, issue, execute, write int
	cu, ib, wf, wg, uop                  int
	asm                                  string
	lifeConcise                          []Activity
	lifeVerbose                          []Activity
}

func (inst *Instruction) sanityCheck(info map[string]string) error {
	// Sanity check
	id, _ := strconv.Atoi(info["id"])
	cu, _ := strconv.Atoi(info["cu"])
	if id != inst.id && cu != inst.cu {
		log.Printf("Expected id/cu=%d/%d, Actual id/cu=%d/%d\n",
			inst.id, inst.cu, id, cu)
		return errors.New("Instruction: id/cu doesn't match!")
	}
	return nil
}

// New record 'New' activity of an instruction
func (inst *Instruction) New(cycle int, info map[string]string) {
	inst.id, _ = strconv.Atoi(info["id"])
	inst.start = cycle
	inst.cu, _ = strconv.Atoi(info["cu"])
	inst.ib, _ = strconv.Atoi(info["ib"])
	inst.wf, _ = strconv.Atoi(info["wf"])
	inst.wg, _ = strconv.Atoi(info["wg"])
	inst.uop, _ = strconv.Atoi(info["uop_id"])
	inst.asm = info["asm"]
	inst.lifeVerbose = append(inst.lifeVerbose, Activity{cycle, info["stg"]})
	inst.lifeConcise = append(inst.lifeConcise, Activity{0, info["stg"]})
}

// Exe record 'Execute' activity of an instruction
func (inst *Instruction) Exe(cycle int, info map[string]string) {
	// Sanity check
	err := inst.sanityCheck(info)
	if err != nil {
		return
	}

	// Get current/previous activity and time interval
	instPrevActivity := inst.lifeVerbose[len(inst.lifeVerbose)-1]
	instCurrActivity := Activity{cycle, info["stg"]}
	timeInterval := instCurrActivity.cycle - instPrevActivity.cycle

	// Record verbose activity
	inst.lifeVerbose = append(inst.lifeVerbose, instCurrActivity)

	// Update concise activity
	inst.lifeConcise[len(inst.lifeConcise)-1].cycle += timeInterval
	if instCurrActivity.activity != instPrevActivity.activity {
		activity := Activity{0, instCurrActivity.activity}
		inst.lifeConcise = append(inst.lifeConcise, activity)
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
	inst.finish = cycle
	inst.length = inst.finish - inst.start
	info["stg"] = "end"
	inst.Exe(cycle, info)

	// Remove the last "end" activity
	inst.lifeConcise = inst.lifeConcise[:len(inst.lifeConcise)-1]
}

// IsValid to check completeness of instruction
func (inst *Instruction) IsValid() bool {
	isValid := true

	// Check field
	isValid = isValid && inst.start != 0
	isValid = isValid && inst.length != 0
	isValid = isValid && inst.asm != ""
	isValid = isValid && len(inst.lifeConcise) != 0
	isValid = isValid && len(inst.lifeVerbose) != 0

	return isValid
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
	instID := "inst_" + strconv.Itoa(inst.id) + "_"
	subgroup := instID

	cycle := inst.start
	for index, activity := range inst.lifeConcise {
		id := instID + activity.activity
		group := inst.cu
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

	id := "inst_" + strconv.Itoa(inst.id) + "_" + strconv.Itoa(inst.cu)
	group := inst.cu
	subgroup := id
	content := inst.asm

	instJSON := InstructionJSON{ID: id, Group: group, Content: content,
		Start: inst.start, End: inst.start + inst.length, SubGroup: subgroup}
	instructionJSONArray = append(instructionJSONArray, &instJSON)

	return instructionJSONArray

}
