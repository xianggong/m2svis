package database

import (
	"strings"

	"github.com/golang/glog"
)

// GetTraceRawdata returns instruction table from database
func GetTraceRawdata(traceName string, filter string) (out []TraceRawdata, err error) {
	useDB()

	insts := []TraceRawdata{}

	query := strings.Join([]string{"SELECT * from", traceName, filter}, " ")
	err = GetInstance().Select(&insts, query)
	if err != nil {
		glog.Warning(err)
		return nil, err
	}

	return insts, err
}
