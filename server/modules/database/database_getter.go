package database

// GetAllTraceInfo returns information of all traces
func GetAllTraceInfo() (out []TraceInfo, err error) {
	useDB()
	data := []TraceInfo{}
	query := "select table_name as TraceName,table_rows as InstCount,create_time as CreateTime from information_schema.Tables where table_schema='m2svis';"
	err = GetInstance().Select(&data, query)

	return data, err
}
