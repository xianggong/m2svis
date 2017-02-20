package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitContainer(t *testing.T) {
	assert := assert.New(t)

	err := initContainer()

	assert.Nil(err)
}
