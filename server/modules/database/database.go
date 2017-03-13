package database

import (
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

const dbName = "m2svis"

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

	// Get database instance
	GetInstance()

	// Connect to database
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
