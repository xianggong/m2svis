package database

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
)

// GetTraceAll returns all traces information
func GetTraceAll() (tracesAll []TraceAll, err error) {
	GetInstance().MustExec("USE m2svis")
	data := []TraceAll{}
	query := "select table_name,table_rows,create_time from information_schema.Tables where table_schema='m2svis';"
	err = GetInstance().Select(&data, query)

	return data, err
}

// GetTraceData returns instruction table from database
func GetTraceData(traceName string, filter string) (out []TraceData, err error) {
	insts := []TraceData{}

	if isDatabaseExist("m2svis") {
		query := "USE m2svis"
		GetInstance().MustExec(query)

		// Get instructions
		query = strings.Join([]string{"SELECT * from", traceName, filter}, " ")
		err = GetInstance().Select(&insts, query)
		if err != nil {
			glog.Warning(err)
			return nil, err
		}
	} else {
		return nil, err
	}

	// Return
	return insts, err
}

// GetTraceCount returns metadata of a trace, such as # of instructions
func GetTraceCount(traceName string, filter string) (out TraceCount, err error) {
	GetInstance().MustExec("USE m2svis")
	traceCount := TraceCount{}
	query := fmt.Sprintf("SELECT count(*) AS CountInsts from %s %s", traceName, filter)
	err = GetInstance().Get(&traceCount, query)
	if err != nil {
		glog.Warning(err)
		return traceCount, err
	}
	return traceCount, err
}

// GetTraceMeta returns metadata of a trace
func GetTraceMeta(traceName, filter string) (out TraceMeta, err error) {
	GetInstance().MustExec("USE m2svis")
	traceCount := TraceMeta{}
	query := fmt.Sprintf("SELECT count(*) as CountInsts, min(st) as MinCycle, max(fn) as MaxCycle, sum(fsw+fsb+ism+isw+isb+dsw+dsb+rsw+rsb+esw+esb+wsw+wsb) as CountStall, count(distinct wf) as CountWF, count(distinct wg) as CountWG, count(distinct cu) as CountCU from %s %s", traceName, filter)
	err = GetInstance().Get(&traceCount, query)
	if err != nil {
		glog.Warning(err)
		return traceCount, err
	}
	return traceCount, err
}

// GetTraceStall returns stall information
func GetTraceStall(traceName, filter string) (out TraceStall, err error) {
	GetInstance().MustExec("USE m2svis")
	traceStall := TraceStall{}
	query := fmt.Sprintf("SELECT sum(fsw+fsb+ism+isw+isb+dsw+dsb+rsw+rsb+esw+esb+wsw+wsb) as StallTotal, sum(fsw+fsb) as StallFrontend, sum(ism+isw+isb) as StallIssue, sum(dsw+dsb) as StallDecode, sum(rsw+rsb) as StallRead, sum(esw+esb) as StallExecute, sum(wsw+wsb) as StallWrite from %s %s", traceName, filter)
	err = GetInstance().Get(&traceStall, query)
	if err != nil {
		glog.Warning(err)
		return traceStall, err
	}
	return traceStall, err
}

// GetTraceNumCU returns number of compute unit
func GetTraceNumCU(traceName string) int {
	GetInstance().MustExec("USE m2svis")
	numCU := 0
	queryCountCU := fmt.Sprintf("select count(distinct cu) from %s", traceName)
	GetInstance().Get(&numCU, queryCountCU)
	return numCU
}

// GetTraceStallRow returns stall information row by row
func GetTraceStallRow(traceName, filter string) ([]TraceStall, error) {
	numCU := GetTraceNumCU(traceName)

	queryStall := []string{}
	for i := 0; i < numCU; i++ {
		queryStall = append(queryStall, fmt.Sprintf("select sum(fsw+fsb+ism+isw+isb+dsw+dsb+rsw+rsb+esw+esb+wsw+wsb) as StallTotal, sum(fsw+fsb) as StallFrontend, sum(ism+isw+isb) as StallIssue, sum(dsw+dsb) as StallDecode, sum(rsw+rsb) as StallRead, sum(esw+esb) as StallExecute, sum(wsw+wsb) as StallWrite from %s where cu = %d", traceName, i))
	}
	query := strings.Join(queryStall, " UNION ALL ")

	result := []TraceStall{}
	err := GetInstance().Select(&result, query)

	return result, err
}

// GetTraceStallColumn returns stall information column by column
func GetTraceStallColumn(traceName, filter string) (TraceStallColumn, error) {
	GetInstance().MustExec("USE m2svis")

	countCU := GetTraceNumCU(traceName)

	queryStallFrontend := []string{}
	queryStallIssue := []string{}
	queryStallDecode := []string{}
	queryStallRead := []string{}
	queryStallExecute := []string{}
	queryStallWrite := []string{}
	for i := 0; i < countCU; i++ {
		queryStallFrontend = append(queryStallFrontend, fmt.Sprintf("select sum(fsw+fsb) as StallFrontend from %s where cu = %d", traceName, i))
		queryStallIssue = append(queryStallIssue, fmt.Sprintf("select sum(ism++isw+isb) as StallIssue from %s where cu = %d", traceName, i))
		queryStallDecode = append(queryStallDecode, fmt.Sprintf("select sum(dsw+dsb) as StallDecode from %s where cu = %d", traceName, i))
		queryStallRead = append(queryStallRead, fmt.Sprintf("select sum(rsw+rsb) as StallRead from %s where cu = %d", traceName, i))
		queryStallExecute = append(queryStallExecute, fmt.Sprintf("select sum(esw+esb) as StallExecute from %s where cu = %d", traceName, i))
		queryStallWrite = append(queryStallWrite, fmt.Sprintf("select sum(wsw+wsb) as StallWrite from %s where cu = %d", traceName, i))
	}
	queryF := strings.Join(queryStallFrontend, " UNION ALL ")
	queryI := strings.Join(queryStallIssue, " UNION ALL ")
	queryD := strings.Join(queryStallDecode, " UNION ALL ")
	queryR := strings.Join(queryStallRead, " UNION ALL ")
	queryE := strings.Join(queryStallExecute, " UNION ALL ")
	queryW := strings.Join(queryStallWrite, " UNION ALL ")

	stallColumn := TraceStallColumn{}
	GetInstance().Select(&stallColumn.Frontend, queryF)
	GetInstance().Select(&stallColumn.Issue, queryI)
	GetInstance().Select(&stallColumn.Decode, queryD)
	GetInstance().Select(&stallColumn.Read, queryR)
	GetInstance().Select(&stallColumn.Execute, queryE)
	GetInstance().Select(&stallColumn.Write, queryW)

	return stallColumn, nil
}

func GetTraceActiveCount(traceName string, cuid, start, finish, windowSize int) ([]int, error) {
	GetInstance().MustExec("USE m2svis")
	st := start
	fn := finish
	meta, _ := GetTraceMeta(traceName, "")
	if start > meta.MaxCycle {
		return nil, nil
	}
	if finish > meta.MaxCycle {
		fn = meta.MaxCycle
	}

	queryActive := []string{}
	for i := st; i < fn; i += windowSize {
		queryActive = append(queryActive, fmt.Sprintf("SELECT COUNT(*) AS CountActiveInstructions FROM %s WHERE st <= %d AND fn >= %d AND cu = %d", traceName, i+st+windowSize, i+st, cuid))
	}
	queryA := strings.Join(queryActive, " UNION ALL ")

	fmt.Println(queryA)

	countActiveInst := []int{}
	GetInstance().Select(&countActiveInst, queryA)
	fmt.Println(countActiveInst)

	return countActiveInst, nil
}
