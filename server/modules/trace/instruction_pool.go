package trace

// InstPool contains instruction objects
type InstPool struct {
	Buffer map[string]*Instruction
}

func (instPool *InstPool) getInst(parseInfo *ParseInfo) *Instruction {
	id := parseInfo.GetID()
	return instPool.Buffer[id]
}

// Process information from parser and process the data, return an instruction
// if it completes all the pipeline stages
func (instPool *InstPool) Process(parseInfo *ParseInfo) (inst *Instruction, err error) {
	// Init insturction buffer if neccesary
	if instPool.Buffer == nil {
		instPool.Buffer = make(map[string]*Instruction)
	}

	// Get cycle and field information
	cycle := parseInfo.cycle
	field := parseInfo.field

	switch parseInfo.key {
	case "si.new_inst":
		// Create new instruction object
		inst := &Instruction{}

		// Update
		err = inst.New(cycle, field)

		// Push to repo
		id := parseInfo.GetID()
		instPool.Buffer[id] = inst

	case "si.inst":
		// Get instruction object
		inst := instPool.getInst(parseInfo)

		// Update
		err = inst.Exe(cycle, field)

	case "si.end_inst":
		// Get instruction object
		inst := instPool.getInst(parseInfo)

		// Update
		err = inst.End(cycle, field)

		// Remove from buffer
		id := parseInfo.GetID()
		delete(instPool.Buffer, id)

		// Return a finished instruction
		return inst, err
	}

	return nil, err
}
