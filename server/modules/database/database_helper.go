package database

import (
	"errors"

	"github.com/golang/glog"
)

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

	// Otherwise it exists
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
