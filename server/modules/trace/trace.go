package trace

import (
	"bufio"
	"compress/gzip"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	// Mysql backend
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

// Trace contains trace related backend modules
type Trace struct {
	Config Configuration

	InstPool InstPool

	Database *sqlx.DB
}

var instance *Trace
var once sync.Once

// GetInstance returns a trace instance in a singleton fashion
func GetInstance() *Trace {
	once.Do(func() {
		instance = &Trace{}
	})
	return instance
}

// Init connects to database backend
func (trace *Trace) Init(configFile string) (err error) {
	// Init configuration and get DSN
	trace.Config.Init(configFile)

	// Connect to database
	trace.Database, err = sqlx.Open("mysql", trace.Config.GetDSN())
	if err != nil {
		glog.Warning(err)
		return err
	}

	// Force a connection and test
	err = trace.Database.Ping()
	if err != nil {
		glog.Warning(err)
		return err
	}

	return nil
}

// Process trace file
func (trace *Trace) Process(path string) (err error) {
	// Get trace file
	file, err := os.Open(path)
	if err != nil {
		glog.Fatal(err)
		return err
	}
	defer file.Close()

	// Open as gzip
	gzfile, err := gzip.NewReader(file)
	if err != nil {
		glog.Fatal(err)
		return err
	}
	defer gzfile.Close()

	// New scanner to read file line by line
	scanner := bufio.NewScanner(gzfile)
	scanner.Split(bufio.ScanLines)

	// Create a database associated with the incoming trace
	trace.NewTraceInDB(path)

	// Query string for inserting instructions to database
	query := "INSERT INTO instructions " + GetInstSQLColNames("", ", ")
	query += " VALUES " + GetInstSQLColNames(":", ",")
	tx := trace.Database.MustBegin()
	parser := new(Parser)
	for scanner.Scan() {
		line := scanner.Text()

		// check for errors
		if err = scanner.Err(); err != nil {
			glog.Fatal(err)
			return err
		}

		// Parse line by line
		info, err := parser.Parse(line)
		if err == nil {
			inst, err := trace.InstPool.Process(&info)
			if inst != nil && err == nil {
				_, err = tx.NamedExec(query, inst)
				if err != nil {
					glog.Error(err)
					return err
				}
			}
		}
	}
	tx.Commit()

	// Make sure no instruction left
	instCount := len(trace.InstPool.Buffer)
	if instCount != 0 {
		glog.Warningf("%s: %d incomplete instruction", path, instCount)
		return errors.New("Some instructions are not processed!")
	}

	// Return
	return nil
}

// NewTraceInDB adds a new database associated with the incoming trace file
func (trace *Trace) NewTraceInDB(path string) {
	// Get filename, remove suffix when neccesary
	fileName := filepath.Base(path)
	fileName = strings.TrimSuffix(fileName, ".gz")

	// Drop database if already exists
	query := "DROP DATABASE IF EXISTS " + fileName
	trace.Database.MustExec(query)

	// Create database with the file name
	query = "CREATE DATABASE IF NOT EXISTS " + fileName
	trace.Database.MustExec(query)

	// Use database
	query = "USE " + fileName
	trace.Database.MustExec(query)

	// Create instruction table in the database
	query = GetSQLQueryNewInstTable("instructions")
	trace.Database.MustExec(query)

}

// GetInstsInDB returns instructions from database
func (trace *Trace) GetInstsInDB(dbName, filter string) (inst []Instruction, err error) {
	insts := []Instruction{}

	// Check if database exists
	query := "SHOW DATABASES LIKE '" + dbName + "'"
	_, err = trace.Database.Queryx(query)
	if err != nil {
		glog.Warning(err)
		return nil, err
	}

	// Use database
	query = "USE " + dbName
	trace.Database.MustExec(query)

	// Get instructions
	query = "SELECT * from instructions " + filter
	err = trace.Database.Select(&insts, query)
	if err != nil {
		glog.Warning(err)
		return nil, err
	}

	// Return
	return insts, err
}
