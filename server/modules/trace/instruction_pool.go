package trace

import "github.com/xianggong/m2svis/server/modules/instruction"

// InstPool contains instruction.Instruction objects
type InstPool struct {
	Buffer map[string]*instruction.Instruction
}

func (instPool *InstPool) getInst(parseInfo *ParseInfo) *instruction.Instruction {
	id := parseInfo.GetID()
	return instPool.Buffer[id]
}

// Process information from parser and process the data, return an instruction.Instruction
// if it completes all the pipeline stages
func (instPool *InstPool) Process(parseInfo *ParseInfo) (inst *instruction.Instruction, err error) {
	// Init insturction buffer if neccesary
	if instPool.Buffer == nil {
		instPool.Buffer = make(map[string]*instruction.Instruction)
	}

	// Get cycle and field information
	cycle := parseInfo.cycle
	field := parseInfo.field

	switch parseInfo.key {
	case "si.new_inst":
		// Create new instruction.Instruction object
		inst := &instruction.Instruction{}

		// Update
		err = inst.New(cycle, field)

		// Push to repo
		id := parseInfo.GetID()
		instPool.Buffer[id] = inst

	case "si.inst":
		// Get instruction.Instruction object
		inst := instPool.getInst(parseInfo)

		// Update
		err = inst.Exe(cycle, field)

	case "si.end_inst":
		// Get instruction.Instruction object
		inst := instPool.getInst(parseInfo)

		// Update
		err = inst.End(cycle, field)

		// Remove from buffer
		id := parseInfo.GetID()
		delete(instPool.Buffer, id)

		// Return a finished instruction.Instruction
		return inst, err
	}

	return nil, err
}
