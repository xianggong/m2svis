package trace

import (
	"bufio"
	"compress/gzip"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/xianggong/m2svis/server/modules/database"
)

// Trace contains instruction pool and parser object
type Trace struct {
	InstPool InstPool
	Parser   Parser
	Config   database.Configuration
	Database *sqlx.DB
}

// Init prepares database
func (trace *Trace) Init(configFile string) (err error) {
	// Init configuration and get DSN
	trace.Config.Init(configFile)

	// Connect to database
	trace.Database, err = sqlx.Open("mysql", trace.Config.GetDSN())
	if err != nil {
		log.Println(err)
		return err
	}

	// Force a connection and test that it worked
	err = trace.Database.Ping()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// Process trace
func (trace *Trace) Process(path string) {
	// Get filename, remove suffix when neccesary
	fileName := filepath.Base(path)
	fileName = strings.TrimSuffix(fileName, ".gz")

	// Create database with the file name
	query := "CREATE DATABASE IF NOT EXISTS " + fileName
	trace.Database.MustExec(query)

	// Use database
	query = "USE " + fileName
	trace.Database.MustExec(query)

	// Create instruction table
	query = GetSQLQueryNewInstTable("instructions")
	trace.Database.MustExec(query)

	// Get trace file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	// Open as gzip
	gzfile, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer gzfile.Close()

	// New scanner to read file line by line
	scanner := bufio.NewScanner(gzfile)
	scanner.Split(bufio.ScanLines)

	// Parse line by line
	for scanner.Scan() {
		line := scanner.Text()

		// check for errors
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}

		info, err := trace.Parser.Parse(line)
		if err == nil {
			inst, err := trace.InstPool.Process(&info)
			if inst != nil && err == nil {
				query = "INSERT INTO instructions " + GetInstructionSQLColumnNames("", ", ")
				query += " VALUES " + GetInstructionSQLColumnNames(":", ",")
				go trace.Database.NamedExec(query, inst)
			}
		}
	}
}

// GetJSON returns JSON arrays
// func (trace *Trace) GetJSON() []*InstructionJSON {
// 	if len(trace.instPool.Ready) == 0 {
// 		return nil
// 	}

// 	var traceJSON []*InstructionJSON

// 	for _, inst := range trace.instPool.Ready {
// 		traceJSON = append(traceJSON, inst.GetOverviewJSON()...)
// 	}

// 	return traceJSON
// }
