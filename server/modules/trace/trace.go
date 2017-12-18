package trace

import (
	"bufio"
	"compress/gzip"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/xianggong/m2svis/server/modules/database"
	"github.com/xianggong/m2svis/server/modules/instruction"
)

// Trace contains trace related backend modules
type Trace struct {
	InstPool InstPool
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

// Init create an instance of trace
func Init() {
	GetInstance()
}

func ProcessGzip(path string) (err error) {
	// Open trace file
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

	// Use m2svis database
	db := database.GetInstance()
	db.MustExec("USE m2svis")

	// Create an instruction table in 'm2svis' database for the incoming trace
	traceName := strings.TrimSuffix(filepath.Base(path), ".gz")
	query := instruction.QueryCreateInstTable(traceName)
	db.MustExec(query)

	// Query string for inserting instructions to database
	query = "INSERT INTO " + traceName + "_insts" + instruction.GetInstSQLColNames("", ", ")
	query += " VALUES " + instruction.GetInstSQLColNames(":", ",")
	tx := db.MustBegin()
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
			inst, err := instance.InstPool.Process(&info)
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
	instCount := len(instance.InstPool.Buffer)
	if instCount != 0 {
		glog.Warningf("%s: %d incomplete instruction", path, instCount)
		return errors.New("instruction incomplete")
	}

	// Return
	return nil
}

// ProcessCSV
func ProcessCSV(path string) (err error) {
	// Use m2svis database
	db := database.GetInstance()
	db.MustExec("USE m2svis")

	// Create an instruction table in 'm2svis' database for the incoming trace
	traceName := strings.TrimSuffix(filepath.Base(path), ".vis")
	query := instruction.GetQueryCreateInstTableCSV(traceName)
	db.MustExec(query)

	// Load data from file to database
	absPath, _ := filepath.Abs(path)
	query = `LOAD DATA LOCAL INFILE "` + absPath + `" INTO TABLE ` + traceName
	query += ` COLUMNS TERMINATED BY '|' `
	query += ` OPTIONALLY ENCLOSED BY '"' `
	query += ` ESCAPED BY '"' `
	query += ` LINES TERMINATED BY '\n'; `
	db.MustExec(query)

	return nil
}

// Process trace file
func Process(path string) (err error) {
	fileExtension := filepath.Ext(path)
	switch fileExtension {
	case ".gz":
		err = ProcessGzip(path)
	case ".vis":
		err = ProcessCSV(path)
	default:
		glog.Errorf("%s extension is not supported\n", fileExtension)
	}

	return err
}
