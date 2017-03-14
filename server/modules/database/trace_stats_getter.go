package database

import (
	"fmt"
	"strings"

	"github.com/golang/glog"
)

// GetTraceCount returns metadata of a trace, such as # of instructions
func GetTraceCount(traceName string, filter string) (out TraceCount, err error) {
	useDB()
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
	useDB()
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
	useDB()
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
	useDB()
	numCU := 0
	queryCountCU := fmt.Sprintf("select count(distinct cu) from %s", traceName)
	GetInstance().Get(&numCU, queryCountCU)
	return numCU
}

// GetTraceStallRow returns stall information row by row
func GetTraceStallRow(traceName, filter string) ([]TraceStall, error) {
	numCU := GetTraceNumCU(traceName)
	useDB()
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
	useDB()

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

// GetTraceActiveCount returns count of active instructions
func GetTraceActiveCount(traceName string, cuid, start, finish, windowSize int) (ActiveInsts, error) {
	useDB()

	st := start
	fn := finish
	meta, _ := GetTraceMeta(traceName, "")
	if start > meta.MaxCycle {
		return ActiveInsts{}, nil
	}
	if finish > meta.MaxCycle {
		fn = meta.MaxCycle
	}

	step := 1
	if windowSize >= 1 {
		step = windowSize
	}

	var activeInsts ActiveInsts
	for i := st; i < fn; i += step {
		count := 0
		queryActive := fmt.Sprintf("SELECT COUNT(*) AS CountActiveInstructions FROM %s WHERE st <= %d AND fn >= %d", traceName, i+st+windowSize, i+st)
		if cuid != -1 {
			queryActive += fmt.Sprintf(" AND cu = %d ", cuid)
		}
		err := GetInstance().Get(&count, queryActive)
		if err != nil {
			glog.Error(err)
		} else {
			activeInsts.Count = append(activeInsts.Count, count)
			activeInsts.Cycle = append(activeInsts.Cycle, i)
		}
	}

	return activeInsts, nil
}

// GetInstCountByInstType returns number of instructions of each type
func GetInstCountByInstType(tracename string) map[string]int {
	useDB()

	queryMap := map[string]string{
		"SOP2":  `SELECT COUNT(*) FROM %s WHERE type="SOP2"`,
		"SOPK":  `SELECT COUNT(*) FROM %s WHERE type="SOPK"`,
		"SOP1":  `SELECT COUNT(*) FROM %s WHERE type="SOP1"`,
		"SOPC":  `SELECT COUNT(*) FROM %s WHERE type="SOPC"`,
		"SOPP":  `SELECT COUNT(*) FROM %s WHERE type="SOPP"`,
		"SMRD":  `SELECT COUNT(*) FROM %s WHERE type="SMRD"`,
		"VOP2":  `SELECT COUNT(*) FROM %s WHERE type="VOP2"`,
		"VOP1":  `SELECT COUNT(*) FROM %s WHERE type="VOP1"`,
		"VOPC":  `SELECT COUNT(*) FROM %s WHERE type="VOPC"`,
		"VOP3A": `SELECT COUNT(*) FROM %s WHERE type="VOP3A"`,
		"VOP3B": `SELECT COUNT(*) FROM %s WHERE type="VOP3B"`,
		"DS":    `SELECT COUNT(*) FROM %s WHERE type="DS"`,
		"MUBUF": `SELECT COUNT(*) FROM %s WHERE type="MUBUF"`,
		"MTBUF": `SELECT COUNT(*) FROM %s WHERE type="MTBUF"`,
		"MIMG":  `SELECT COUNT(*) FROM %s WHERE type="MIMG"`,
	}

	// Get from database
	instCountByType := map[string]int{}
	for key, value := range queryMap {
		count := 0
		err := GetInstance().Get(&count, fmt.Sprintf(value, tracename))
		if err != nil {
			glog.Error(err)
			return map[string]int{}
		}
		instCountByType[key] = count
	}

	return instCountByType
}

// GetInstCountByExecUnit returns number of instructions on each execution unit
func GetInstCountByExecUnit(tracename string) map[string]int {
	useDB()

	queryMap := map[string]string{
		"Branch":       `SELECT COUNT(*) FROM %s WHERE eu="Branch"`,
		"LDS":          `SELECT COUNT(*) FROM %s WHERE eu="LDS"`,
		"VectorMemory": `SELECT COUNT(*) FROM %s WHERE eu="VectorMemory"`,
		"Scalar":       `SELECT COUNT(*) FROM %s WHERE eu="Scalar"`,
		"SIMD":         `SELECT COUNT(*) FROM %s WHERE eu="SIMD"`,
	}

	// Get from database
	instCountByExecUnit := map[string]int{}
	for key, value := range queryMap {
		count := 0
		err := GetInstance().Get(&count, fmt.Sprintf(value, tracename))
		if err != nil {
			glog.Error(err)
			return map[string]int{}
		}
		instCountByExecUnit[key] = count
	}

	return instCountByExecUnit
}

// GetCycleCountByInstType returns number of instructions on each execution unit
func GetCycleCountByInstType(tracename string) map[string]int {
	useDB()

	queryMap := map[string]string{
		"SOP2":  `SELECT COALESCE(SUM(len), 0) FROM %s WHERE type="SOP2"`,
		"SOPK":  `SELECT COALESCE(SUM(len), 0) FROM %s WHERE type="SOPK"`,
		"SOP1":  `SELECT COALESCE(SUM(len), 0) FROM %s WHERE type="SOP1"`,
		"SOPC":  `SELECT COALESCE(SUM(len), 0) FROM %s WHERE type="SOPC"`,
		"SOPP":  `SELECT COALESCE(SUM(len), 0) FROM %s WHERE type="SOPP"`,
		"SMRD":  `SELECT COALESCE(SUM(len), 0) FROM %s WHERE type="SMRD"`,
		"VOP2":  `SELECT COALESCE(SUM(len), 0) FROM %s WHERE type="VOP2"`,
		"VOP1":  `SELECT COALESCE(SUM(len), 0) FROM %s WHERE type="VOP1"`,
		"VOPC":  `SELECT COALESCE(SUM(len), 0) FROM %s WHERE type="VOPC"`,
		"VOP3A": `SELECT COALESCE(SUM(len), 0) FROM %s WHERE type="VOP3A"`,
		"VOP3B": `SELECT COALESCE(SUM(len), 0) FROM %s WHERE type="VOP3B"`,
		"DS":    `SELECT COALESCE(SUM(len), 0) FROM %s WHERE type="DS"`,
		"MUBUF": `SELECT COALESCE(SUM(len), 0) FROM %s WHERE type="MUBUF"`,
		"MTBUF": `SELECT COALESCE(SUM(len), 0) FROM %s WHERE type="MTBUF"`,
		"MIMG":  `SELECT COALESCE(SUM(len), 0) FROM %s WHERE type="MIMG"`,
	}

	// Get from database
	cycleCountByType := map[string]int{}
	for key, value := range queryMap {
		count := 0
		err := GetInstance().Get(&count, fmt.Sprintf(value, tracename))
		if err != nil {
			glog.Error(err)
			return map[string]int{}
		}
		cycleCountByType[key] = count
	}

	return cycleCountByType
}

// GetCycleCountByExecUnit returns number of instructions on each execution unit
func GetCycleCountByExecUnit(tracename string) map[string]int {
	useDB()

	queryMap := map[string]string{
		"Branch":       `SELECT COALESCE(SUM(len), 0) FROM %s WHERE eu="Branch"`,
		"LDS":          `SELECT COALESCE(SUM(len), 0) FROM %s WHERE eu="LDS"`,
		"VectorMemory": `SELECT COALESCE(SUM(len), 0) FROM %s WHERE eu="VectorMemory"`,
		"Scalar":       `SELECT COALESCE(SUM(len), 0) FROM %s WHERE eu="Scalar"`,
		"SIMD":         `SELECT COALESCE(SUM(len), 0) FROM %s WHERE eu="SIMD"`,
	}

	// Get from database
	cycleCountByExecUnit := map[string]int{}
	for key, value := range queryMap {
		count := 0
		query := fmt.Sprintf(value, tracename)
		err := GetInstance().Get(&count, query)
		if err != nil {
			glog.Error(err)
			return map[string]int{}
		}
		cycleCountByExecUnit[key] = count
	}

	return cycleCountByExecUnit
}

func GetCycleCountByCU(tracename string) map[int]int {
	useDB()

	cycleCountByCU := map[int]int{}
	countCU := GetTraceNumCU(tracename)
	for cuid := 0; cuid < countCU; cuid++ {
		len := 0
		query := fmt.Sprintf("SELECT MAX(fn) FROM %s WHERE cu=%d", tracename, cuid)
		err := GetInstance().Get(&len, query)
		if err != nil {
			glog.Error(err)
			return map[int]int{}
		}
		cycleCountByCU[cuid] = len
	}

	return cycleCountByCU
}
