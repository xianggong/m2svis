package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraceInit(t *testing.T) {
	assert := assert.New(t)
	trace := GetInstance()

	err := trace.Init("test/config.toml")

	assert.Nil(err)
}

func TestTraceProcess(t *testing.T) {
	assert := assert.New(t)
	trace := GetInstance()

	err := trace.Init("test/config.toml")
	assert.Nil(err)

	err = trace.Process("test/test.gz")
	assert.Nil(err)
}

func TestTraceGetInsts(t *testing.T) {
	assert := assert.New(t)
	trace := GetInstance()

	err := trace.Init("test/config.toml")
	assert.Nil(err)

	insts, err := trace.GetInstsInDB("test", "")
	assert.Nil(err)
	assert.NotEqual(0, len(insts))
}
