package database

import "testing"
import "github.com/stretchr/testify/assert"

func TestInit(t *testing.T) {
	assert := assert.New(t)

	var config Configuration
	config.Init("test/config.toml")

	assert.Equal("root", config.Username)
	assert.Equal("mysqltest", config.Password)
	assert.Equal("tcp", config.Protocol)
	assert.Equal("127.0.0.1", config.Address)
	assert.Equal("3306", config.Port)
}

func TestGetDSN(t *testing.T) {
	assert := assert.New(t)

	var config Configuration
	config.Init("test/config.toml")
	dsn := config.GetDSN()

	assert.Equal("root:mysqltest@tcp(127.0.0.1:3306)/", dsn)
}
