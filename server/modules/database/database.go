package database

import (
	"database/sql"
	"log"

	// Mysql backend
	_ "github.com/go-sql-driver/mysql"
)

// Database backend for storing and retriving data
type Database struct {
	config   Configuration
	database *sql.DB
}

// Open opens a database and returns it
func (db *Database) Open(dbName ...string) (*sql.DB, error) {
	var err error

	// Get data source name
	dsn := db.config.GetDSN()
	if len(dbName) == 1 {
		dsn += dbName[0]
	}

	// Open database
	db.database, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Return
	return db.database, err
}

// Close internal database if it is still open
func (db *Database) Close() error {
	if db.database != nil {
		return db.database.Close()
	}

	return nil
}

// NewDB creates a new database with the name
func (db *Database) NewDB(dbname string) (sql.Result, error) {
	query := "CREATE DATABASE IF NOT EXISTS " + dbname + ";"

	return db.database.Exec(query)
}

// UseDB use database
func (db *Database) UseDB(dbname string) (sql.Result, error) {
	query := "USE " + dbname + ";"

	return db.database.Exec(query)
}

// DelDB delete/drop a database
func (db *Database) DelDB(dbname string) (sql.Result, error) {
	query := "DROP DATABASE " + dbname + ";"

	return db.database.Exec(query)
}

// NewInstTable creates an instruction table
func (db *Database) NewInstTable(instTableName string) (sql.Result, error) {
	query := "CREATE TABLE " + instTableName + "("
	query += "id INTEGER,"
	query += "start INTEGER,"
	query += "end INTEGER,"
	query += "length INTEGER,"
	query += "cu INTEGER,"
	query += "ib INTEGER,"
	query += "wf INTEGER,"
	query += "wg INTEGER,"
	query += "uop INTEGER"
	query += ");"

	return db.database.Exec(query)
}

// DelInstTable drops/deletes an instruction table
func (db *Database) DelInstTable(instTableName string) (sql.Result, error) {
	query := "DROP TABLE " + instTableName + ";"

	return db.database.Exec(query)
}

// InsertInstTable insert values into an instruction table
func (db *Database) InsertInstTable(instTableName string, args ...interface{}) (sql.Result, error) {
	query := "INSERT INTO " + instTableName + " VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);"

	return db.database.Exec(query, args...)
}

// QueryInstTable returns query of an instruction table
func (db *Database) QueryInstTable(instTableName string, conditions string, args ...interface{}) (*sql.Rows, error) {
	query := "SELECT * FROM " + instTableName
	if conditions != "" {
		query += "WHERE " + conditions
	}

	return db.database.Query(query, args...)
}
