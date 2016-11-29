package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpen(t *testing.T) {
	assert := assert.New(t)

	var db Database
	db.config.Init("test/config.toml")
	db.Open()
	assert.NotNil(db.database)
}

func TestNewDatabase(t *testing.T) {
	assert := assert.New(t)

	var db Database
	db.config.Init("test/config.toml")
	database, _ := db.Open()
	defer database.Close()

	_, err := db.NewDB("m2svistest")
	assert.Nil(err)
}

func TestUseDatabase(t *testing.T) {
	assert := assert.New(t)

	var db Database
	db.config.Init("test/config.toml")
	database, _ := db.Open()
	defer database.Close()

	_, err := db.UseDB("m2svistest")
	assert.Nil(err)
}

func TestDelDatabase(t *testing.T) {
	assert := assert.New(t)

	var db Database
	db.config.Init("test/config.toml")
	database, _ := db.Open()
	defer database.Close()

	_, err := db.DelDB("m2svistest")
	assert.Nil(err)
}

func TestNewInstTable(t *testing.T) {
	assert := assert.New(t)

	var db Database
	db.config.Init("test/config.toml")
	database, _ := db.Open()
	defer database.Close()

	_, err := db.NewDB("m2svistest")
	assert.Nil(err)

	_, err = db.UseDB("m2svistest")
	assert.Nil(err)

	_, err = db.NewInstTable("testInst")
	assert.Nil(err)
}

func TestInsertInstTable(t *testing.T) {
	assert := assert.New(t)

	var db Database
	db.config.Init("test/config.toml")
	database, _ := db.Open()
	defer database.Close()

	_, err := db.UseDB("m2svistest")
	assert.Nil(err)

	_, err = db.InsertInstTable("testInst", 0, 1, 2, 3, 4, 5, 6, 7, 8)
	assert.Nil(err)
}

func TestQueryInstTable(t *testing.T) {
	assert := assert.New(t)

	var db Database
	db.config.Init("test/config.toml")
	database, _ := db.Open()
	defer database.Close()

	_, err := db.UseDB("m2svistest")
	assert.Nil(err)

	_, err = db.InsertInstTable("testInst", 0, 1, 2, 3, 4, 5, 6, 7, 8)
	assert.Nil(err)

	rows, err := db.QueryInstTable("testInst", "")
	assert.Nil(err)

	for rows.Next() {
		var id, start, end, length, cu, ib, wf, wg, uop int
		err = rows.Scan(&id, &start, &end, &length, &cu, &ib, &wf, &wg, &uop)
		assert.Equal(0, id)
		assert.Equal(1, start)
		assert.Equal(2, end)
		assert.Equal(3, length)
		assert.Equal(4, cu)
		assert.Equal(5, ib)
		assert.Equal(6, wf)
		assert.Equal(7, wg)
		assert.Equal(8, uop)
	}
	assert.Nil(err)

}

func TestDelInstTable(t *testing.T) {
	assert := assert.New(t)

	var db Database
	db.config.Init("test/config.toml")
	database, _ := db.Open()
	defer database.Close()

	_, err := db.UseDB("m2svistest")
	assert.Nil(err)

	_, err = db.DelInstTable("testInst")
	assert.Nil(err)
}
