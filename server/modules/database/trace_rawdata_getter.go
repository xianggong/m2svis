package database

import (
	"strings"

	"github.com/golang/glog"
	"github.com/xianggong/m2svis/server/modules/instruction"
)

// GetTraceRawdata returns instruction table from database
func GetTraceRawdata(traceName string, filter string) (out []instruction.InstructionCSV, err error) {
	insts := []instruction.InstructionCSV{}

	query := strings.Join([]string{"SELECT * from", traceName, filter}, " ")
	err = GetInstance().Select(&insts, query)
	if err != nil {
		glog.Warning(err)
		return nil, err
	}

	return insts, err
}
