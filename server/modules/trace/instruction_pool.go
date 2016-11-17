package trace

// InstructionPool contains instruction objects
type InstructionPool struct {
	Ready      []*Instruction
	InProgress map[string]*Instruction
}

var instFuncTable = map[string]func(*Instruction, int, map[string]string){
	"si.new_inst": (*Instruction).New,
	"si.inst":     (*Instruction).Exe,
	"si.end_inst": (*Instruction).End,
}

func (instPool *InstructionPool) getInst(parseInfo *ParseInfo) *Instruction {
	id := parseInfo.GetID()
	return instPool.InProgress[id]
}

// Receive information from parser and process the data
func (instPool *InstructionPool) Receive(parseInfo *ParseInfo) (err error) {
	if instPool.InProgress == nil {
		instPool.InProgress = make(map[string]*Instruction)
	}

	if instPool.Ready == nil {
		instPool.Ready = []*Instruction{}
	}

	// Get cycle and field information
	cycle := parseInfo.cycle
	field := parseInfo.field

	switch parseInfo.key {
	case "si.new_inst":
		// Create new instruction object
		inst := &Instruction{}

		// Update
		inst.New(cycle, field)

		// Push to repo
		id := parseInfo.GetID()
		instPool.InProgress[id] = inst

	case "si.inst":
		// Get instruction object
		inst := instPool.getInst(parseInfo)

		// Update
		inst.Exe(cycle, field)

	case "si.end_inst":
		// Get instruction object
		inst := instPool.getInst(parseInfo)

		// Update
		inst.End(cycle, field)

		// fmt.Printf("%d %d [%d-%d]: %s\n", inst.id, inst.cu, inst.start,
		// 	inst.length, inst.asm)
		// fmt.Println(inst.lifeConcise)

		// Push to ready pool and remove from progress pool
		id := parseInfo.GetID()
		instPool.Ready = append(instPool.Ready, inst)
		delete(instPool.InProgress, id)
	}

	return nil
}
