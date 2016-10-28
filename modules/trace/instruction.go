package trace

import "strconv"

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
}

// Exe record 'Execute' activity of an instruction
func (inst *Instruction) Exe(cycle int, info map[string]string) {
	instPrevActivity := inst.lifeVerbose[len(inst.lifeVerbose)-1]
	instCurrActivity := Activity{cycle, info["stg"]}

	// Record verbose activity
	inst.lifeVerbose = append(inst.lifeVerbose, instCurrActivity)

	// Update concise activity
	if instCurrActivity.activity != instPrevActivity.activity {
		activityLength := instCurrActivity.cycle - instPrevActivity.cycle
		activity := Activity{activityLength, instPrevActivity.activity}
		inst.lifeConcise = append(inst.lifeConcise, activity)
	}
}

// End record 'End' activity of an instruction
func (inst *Instruction) End(cycle int, info map[string]string) {
	inst.finish = cycle
	inst.length = inst.finish - inst.start
	inst.Exe(cycle, info)
}
