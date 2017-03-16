package database

import "testing"
import "github.com/stretchr/testify/assert"

func TestRead(t *testing.T) {
	assert := assert.New(t)

	var conf config
	conf.read("../../config.json")

	assert.Equal("m2svis", conf.Container)
	assert.Equal("m2svis", conf.Username)
	assert.Equal("m2svis", conf.Password)
	assert.Equal("m2svis", conf.Database)
	assert.Equal("tcp", conf.Protocol)
	assert.Equal("127.0.0.1", conf.Address)
	assert.Equal("3306", conf.Port)
}

func TestGetDSN(t *testing.T) {
	assert := assert.New(t)

	var conf config
	conf.read("../../config.json")
	dsn := conf.getDSN()

	assert.Equal("m2svis:m2svis@tcp(127.0.0.1:3306)/m2svis?parseTime=true", dsn)
}
