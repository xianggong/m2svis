package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInstance(t *testing.T) {
	assert := assert.New(t)

	db := GetInstance()

	assert.NotNil(db)
}

func TestInit(t *testing.T) {
	assert := assert.New(t)

	err := Init("../../config.json")
	assert.Nil(err, "Init")
}

func TestIsDatabaseExist(t *testing.T) {
	assert := assert.New(t)

	isExist := isDatabaseExist("m2svis")
	assert.Equal(true, isExist)

	isExist = isDatabaseExist("invalid")
	assert.Equal(false, isExist)
}

func TestIsTableExist(t *testing.T) {
	assert := assert.New(t)

	err := isTableExist("m2svis", "noexisttable")
	assert.NotNil(err, "isTableExist")

	err = isTableExist("invalid", "noexisttable")
	assert.NotNil(err, "isTableExist")
}

// func TestGetInstTable(t *testing.T) {
// 	assert := assert.New(t)
//
// 	data, err := GetInstTable("test", "")
// 	assert.NotNil(data, "Received instructions")
// 	assert.Nil(err, "")
// }
