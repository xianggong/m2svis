package database

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

// GetTraceCount returns metadata of a trace, such as # of instructions
func GetTraceCount(traceName string, filter string) (out TraceCount, err error) {
	traceCount := TraceCount{}
	query := fmt.Sprintf("SELECT count(*) AS CountInsts from %s %s", traceName, filter)
	err = GetInstance().Get(&traceCount, query)
	if err != nil {
		glog.Warning(err)
	}
	return traceCount, err
}

// GetTraceMeta returns metadata of a trace
func GetTraceMeta(traceName, filter string) (out TraceMeta, err error) {
	traceCount := TraceMeta{}
	query := fmt.Sprintf("SELECT count(*) as CountInsts, min(Start) as MinCycle, max(Finish) as MaxCycle, sum(FetchStall+IssueStall+DecodeStall+ReadStall+ExecuteStall+WriteStall) as CountStall, count(distinct WF) as CountWF, count(distinct WG) as CountWG, count(distinct CU) as CountCU from %s %s", traceName, filter)
	err = GetInstance().Get(&traceCount, query)
	if err != nil {
		glog.Warning(err)
		return traceCount, err
	}
	return traceCount, err
}

// GetTraceStall returns stall information
func GetTraceStall(traceName, filter string) (out TraceStall, err error) {
	traceStall := TraceStall{}
	query := fmt.Sprintf("SELECT sum(FetchStall+IssueStall+DecodeStall+ReadStall+ExecuteStall+WriteStall) as StallTotal, sum(FetchStall) as StallFrontend, sum(IssueStall) as StallIssue, sum(DecodeStall) as StallDecode, sum(ReadStall) as StallRead, sum(ExecuteStall) as StallExecute, sum(WriteStall) as StallWrite from %s %s", traceName, filter)
	err = GetInstance().Get(&traceStall, query)
	if err != nil {
		glog.Warning(err)
		return traceStall, err
	}
	return traceStall, err
}

// GetTraceNumCU returns number of compute unit
func GetTraceNumCU(traceName string) int {
	numCU := 0
	queryCountCU := fmt.Sprintf("SELECT count(distinct CU) from %s", traceName)
	GetInstance().Get(&numCU, queryCountCU)
	return numCU
}

// GetTraceStallColumn returns stall information column by column
func GetTraceStallColumn(traceName, filter string) (TraceStallColumn, error) {
	queryF := fmt.Sprintf("select sum(FetchStall) as StallFrontend from %s GROUP BY cu", traceName)
	queryI := fmt.Sprintf("select sum(IssueStall) as StallIssue from %s GROUP BY cu", traceName)
	queryD := fmt.Sprintf("select sum(DecodeStall) as StallDecode from %s GROUP BY cu", traceName)
	queryR := fmt.Sprintf("select sum(ReadStall) as StallRead from %s GROUP BY cu", traceName)
	queryE := fmt.Sprintf("select sum(ExecuteStall) as StallExecute from %s GROUP BY cu", traceName)
	queryW := fmt.Sprintf("select sum(WriteStall) as StallWrite from %s GROUP BY cu", traceName)

	stallColumn := TraceStallColumn{}
	GetInstance().Select(&stallColumn.Frontend, queryF)
	GetInstance().Select(&stallColumn.Issue, queryI)
	GetInstance().Select(&stallColumn.Decode, queryD)
	GetInstance().Select(&stallColumn.Read, queryR)
	GetInstance().Select(&stallColumn.Execute, queryE)
	GetInstance().Select(&stallColumn.Write, queryW)

	return stallColumn, nil
}

// GetInstCountByInstType returns number of instructions of each type
func GetInstCountByInstType(tracename string) map[string]int {
	instType := GetInstType(tracename)
	template := `SELECT COUNT(*) FROM %s WHERE Type = "%s"`

	// Get from database
	instCountByType := map[string]int{}
	for _, instType := range instType {
		count := 0
		err := GetInstance().Get(&count, fmt.Sprintf(template, tracename, instType))
		if err != nil {
			glog.Error(err)
			return map[string]int{}
		}
		instCountByType[instType] = count
	}

	return instCountByType
}

// GetInstCountByExecUnit returns number of instructions on each execution unit
func GetInstCountByExecUnit(tracename string) map[string]int {
	queryMap := map[string]string{
		"Branch":       `SELECT COUNT(*) FROM %s WHERE ExecutionUnit="branch"`,
		"LDS":          `SELECT COUNT(*) FROM %s WHERE ExecutionUnit="lds"`,
		"VectorMemory": `SELECT COUNT(*) FROM %s WHERE ExecutionUnit="simd-m"`,
		"Scalar":       `SELECT COUNT(*) FROM %s WHERE ExecutionUnit="scalar"  `,
		"SIMD":         `SELECT COUNT(*) FROM %s WHERE ExecutionUnit="simd"`,
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

// GetInstLengthByExecUnit returns number of instructions on each execution unit
func GetInstLengthByExecUnit(tracename string) map[string]float64 {
	queryMap := map[string]string{
		"Branch":       `SELECT AVG(Length) FROM %s WHERE ExecutionUnit="branch"`,
		"LDS":          `SELECT AVG(Length) FROM %s WHERE ExecutionUnit="lds"`,
		"VectorMemory": `SELECT AVG(Length) FROM %s WHERE ExecutionUnit="simd-m"`,
		"Scalar":       `SELECT AVG(Length) FROM %s WHERE ExecutionUnit="scalar"  `,
		"SIMD":         `SELECT AVG(Length) FROM %s WHERE ExecutionUnit="simd"`,
	}

	// Get from database
	instLengthByExecUnit := map[string]float64{}
	for key, value := range queryMap {
		length := []uint8{}
		err := GetInstance().Get(&length, fmt.Sprintf(value, tracename))
		if err != nil {
			glog.Error(err)
			return map[string]float64{}
		}
		lengthFloat, _ := strconv.ParseFloat(string(length), 64)
		instLengthByExecUnit[key] = lengthFloat
	}

	return instLengthByExecUnit
}

func GetInstType(tracename string) []string {
	template := "SELECT Type from %s GROUP BY Type"
	query := fmt.Sprintf(template, tracename)
	result := []string{}
	err := GetInstance().Select(&result, query)
	if err != nil {
		glog.Error(err)
		return []string{}
	}
	return result
}

func GetExecUnit(tracename string) []string {
	template := "SELECT ExecutionUnit from %s GROUP BY ExecutionUnit"
	query := fmt.Sprintf(template, tracename)
	result := []string{}
	err := GetInstance().Select(&result, query)
	if err != nil {
		glog.Error(err)
		return []string{}
	}
	return result
}

// GetInstLengthByInstType returns number of instructions on each execution unit
func GetInstLengthByInstType(tracename string) map[string][]uint {
	instType := GetInstType(tracename)
	evalFunc := []string{"Min", "Max", "Avg"}
	template := `SELECT CAST(%s(Length) AS UNSIGNED) AS Length FROM %s WHERE Type="%s"`

	result := map[string][]uint{}
	for _, evalFunc := range evalFunc {
		info := []uint{}
		queries := []string{}
		for _, instType := range instType {
			queries = append(queries, fmt.Sprintf(template, evalFunc, tracename, instType))
		}
		query := strings.Join(queries, " UNION ALL ")
		err := GetInstance().Select(&info, query)
		if err != nil {
			glog.Error(err)
		}
		result[evalFunc] = info
	}

	return result
}

// GetCycleCountByInstType returns number of instructions on each execution unit
func GetCycleCountByInstType(tracename string) map[string]int {
	instType := GetInstType(tracename)
	template := `SELECT COALESCE(SUM(Length), 0) FROM %s WHERE Type = "%s"`

	// Get from database
	cycleCountByType := map[string]int{}
	for _, instType := range instType {
		count := 0
		err := GetInstance().Get(&count, fmt.Sprintf(template, tracename, instType))
		if err != nil {
			glog.Error(err)
			return map[string]int{}
		}
		cycleCountByType[instType] = count
	}

	return cycleCountByType
}

// GetCycleCountByExecUnit returns number of instructions on each execution unit
func GetCycleCountByExecUnit(tracename string) map[string]int {
	queryMap := map[string]string{
		"Branch":       `SELECT COALESCE(SUM(Length), 0) FROM %s WHERE ExecutionUnit="branch"`,
		"LDS":          `SELECT COALESCE(SUM(Length), 0) FROM %s WHERE ExecutionUnit="lds"`,
		"VectorMemory": `SELECT COALESCE(SUM(Length), 0) FROM %s WHERE ExecutionUnit="simd-m"`,
		"Scalar":       `SELECT COALESCE(SUM(Length), 0) FROM %s WHERE ExecutionUnit="scalar"`,
		"SIMD":         `SELECT COALESCE(SUM(Length), 0) FROM %s WHERE ExecutionUnit="simd"`,
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

// GetCycleCountByCU returns exectution time
func GetCycleCountByCU(tracename string) map[int]int {
	cycleCountByCU := map[int]int{}
	countCU := GetTraceNumCU(tracename)
	for cuid := 0; cuid < countCU; cuid++ {
		len := 0
		query := fmt.Sprintf("SELECT MAX(Finish) - MIN(Start) FROM %s WHERE CU=%d", tracename, cuid)
		err := GetInstance().Get(&len, query)
		if err != nil {
			glog.Error(err)
			return map[int]int{}
		}
		cycleCountByCU[cuid] = len
	}

	return cycleCountByCU
}
