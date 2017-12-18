package database

import (
	"fmt"

	"github.com/golang/glog"
)

// GetTraceActiveCount returns count of active instructions
func GetTraceActiveCount(traceName string, cuid, start, finish, windowSize int) (map[string][]int, error) {
	st := start
	fn := finish
	meta, _ := GetTraceMeta(traceName, "")
	if start > meta.MaxCycle {
		return make(map[string][]int), nil
	}
	if start < meta.MinCycle {
		st = meta.MinCycle
	}
	if finish > meta.MaxCycle {
		fn = meta.MaxCycle
	}

	step := 1
	if windowSize >= 1 {
		step = windowSize
	}

	queryMap := map[string]string{
		"Branch":       `SELECT COUNT(*) FROM %s WHERE Start <= %d AND Finish >= %d AND ExecutionUnit="branch" `,
		"LDS":          `SELECT COUNT(*) FROM %s WHERE Start <= %d AND Finish >= %d AND ExecutionUnit="lds" `,
		"VectorMemory": `SELECT COUNT(*) FROM %s WHERE Start <= %d AND Finish >= %d AND ExecutionUnit="simd-m" `,
		"Scalar":       `SELECT COUNT(*) FROM %s WHERE Start <= %d AND Finish >= %d AND ExecutionUnit="scalar" `,
		"SIMD":         `SELECT COUNT(*) FROM %s WHERE Start <= %d AND Finish >= %d AND ExecutionUnit="simd" `,
	}

	activeInstsMap := make(map[string][]int)

	// Cycle information
	activeCycle := []int{}
	for i := st; i <= fn; i += step {
		activeCycle = append(activeCycle, i)
	}
	activeInstsMap["Cycle"] = activeCycle

	// Count information
	for key, value := range queryMap {
		activeInsts := []int{}
		for i := 0; i <= fn-st; i += step {
			count := 0
			query := fmt.Sprintf(value, traceName, i+st+windowSize, i+st)
			if cuid != -1 {
				query += fmt.Sprintf(" AND CU = %d ", cuid)
			}
			err := GetInstance().Get(&count, query)
			if err != nil {
				glog.Error(err)
			} else {
				activeInsts = append(activeInsts, count)
			}
		}
		activeInstsMap[key] = activeInsts
	}

	return activeInstsMap, nil
}

func GetExecUnitUtilization(tracename string, cuid, start, finish, windowSize int) map[string]float32 {
	st := start
	fn := finish
	meta, _ := GetTraceMeta(tracename, "")
	if st > meta.MaxCycle {
		st = meta.MaxCycle
	}
	if fn > meta.MaxCycle {
		fn = meta.MaxCycle
	}

	step := 1
	if windowSize >= 1 {
		step = windowSize
	}

	activeCount := map[string]int{
		"branch": 0,
		"lds":    0,
		"simd-m": 0,
		"scalar": 0,
		"simd":   0,
	}

	for i := 0; i <= fn-st; i += step {
		count := []string{}
		query := fmt.Sprintf("SELECT DISTINCT ExecutionUnit FROM %s WHERE Start <= %d AND Finish >= %d", tracename, i+st+windowSize, i+st)
		if cuid != -1 {
			query += fmt.Sprintf(" AND cu = %d ", cuid)
		}
		err := GetInstance().Select(&count, query)
		if err != nil {
			glog.Error(err)
		} else {
			for _, val := range count {
				activeCount[val]++
			}
		}
	}

	length := float32(finish - start)
	utilization := map[string]float32{
		"branch": float32(activeCount["branch"]) / length,
		"lds":    float32(activeCount["lds"]) / length,
		"simd-m": float32(activeCount["simd-m"]) / length,
		"scalar": float32(activeCount["scalar"]) / length,
		"simd":   float32(activeCount["simd"]) / length,
	}

	return utilization
}
