package trace

import (
	"github.com/jmoiron/sqlx"
	"github.com/xianggong/m2svis/server/modules/database"
)

// InstPool contains instruction objects
type InstPool struct {
	Buffer   map[string]*Instruction
	Config   database.Configuration
	Database *sqlx.DB
}

// Init connect to database specified in the configuration file
func (instPool *InstPool) Init(configFile string) {
	// Get data source name
	dsn := instPool.Config.GetDSN()

	// Connect to database
	instPool.Database, _ = sqlx.Connect("mysql", dsn)

	// Create table in database
	inst := &Instruction{}
	instPool.Database.MustExec(inst.GetSQLQueryInsertTable("test"))
}

func (instPool *InstPool) getInst(parseInfo *ParseInfo) *Instruction {
	id := parseInfo.GetID()
	return instPool.Buffer[id]
}

// Process information from parser and process the data
func (instPool *InstPool) Process(parseInfo *ParseInfo) (err error) {
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
		inst.New(cycle, field)

		// Push to repo
		id := parseInfo.GetID()
		instPool.Buffer[id] = inst

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

		// Save to database and remove from buffer
		id := parseInfo.GetID()
		delete(instPool.Buffer, id)
	}

	return nil
}
