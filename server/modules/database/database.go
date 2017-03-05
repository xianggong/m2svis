package database

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
	// Mysql backend
	_ "github.com/go-sql-driver/mysql"
)

// Database
var dbConfig config
var dbInstance *sqlx.DB
var once sync.Once

// GetInstance returns an instance of database in a singleton fashion
func GetInstance() *sqlx.DB {
	once.Do(func() {
		dbInstance = &sqlx.DB{}
	})
	return dbInstance
}

// Init connects to database backend
func Init(configFile string) (err error) {
	// Init configuration
	dbConfig.read(configFile)

	// Connect to database
	GetInstance()
	dbInstance, err = sqlx.Open("mysql", dbConfig.getDSN())
	if err != nil {
		glog.Fatal(err)
		return err
	}

	// Force a connection and test
	err = dbInstance.Ping()
	if err != nil {
		glog.Fatal(err)
		return err
	}

	return nil
}

// Check if database exists
func isDatabaseExist(dbName string) (isExist bool) {
	query := "SHOW DATABASES LIKE '" + dbName + "'"
	rows, err := GetInstance().Queryx(query)

	// Error means does not exists
	if err != nil {
		glog.Warning(err)
		return false
	}

	// No row means does not exists
	if !rows.Next() {
		return false
	}
	return true
}

func isTableExist(dbName, tbName string) (err error) {
	// Check if database exists
	isExist := isDatabaseExist(dbName)
	if !isExist {
		err = errors.New("database does not exist")
		glog.Warning(err)
		return err
	}

	// Database exists, now check table
	query := "USE " + dbName
	GetInstance().MustExec(query)

	count := 0
	query = "SELECT COUNT(*) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_NAME = N'" + tbName + "'"
	GetInstance().Select(&count, query)
	if err != nil {
		glog.Warning(err)
		return err
	}
	if count == 0 {
		return errors.New("no table")
	}

	// No error
	return nil
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
	query := fmt.Sprintf("SELECT count(*) from %s %s", traceName, filter)
	err = GetInstance().Get(&traceCount, query)
	if err != nil {
		glog.Warning(err)
		return traceCount, err
	}

	return traceCount, err
}

// GetTraceAll returns all traces information
func GetTraceAll() (tracesAll []TraceAll, err error) {
	GetInstance().MustExec("USE m2svis")
	data := []TraceAll{}
	query := "select table_name,table_rows,create_time from information_schema.Tables where table_schema='m2svis';"
	err = GetInstance().Select(&data, query)

	return data, err
}
