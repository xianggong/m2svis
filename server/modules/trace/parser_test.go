package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseClock(t *testing.T) {
	assert := assert.New(t)
	parser := new(Parser)

	input := "c clk=1000"
	actual, err := parser.Parse(input)

	assert.Equal(nil, err)
	assert.Equal("1000", actual.field["clock"])
}

func TestParseInstructionNew(t *testing.T) {
	assert := assert.New(t)
	parser := new(Parser)

	input := `si.new_inst id=1 cu=2 ib=3 wg=4 wf=5 uop_id=6 stg="f" ` +
		`asm="s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358"`
	actual, err := parser.Parse(input)

	assert.Equal(nil, err)
	assert.Equal("si", actual.field["arch"])
	assert.Equal("1", actual.field["id"])
	assert.Equal("2", actual.field["cu"])
	assert.Equal("3", actual.field["ib"])
	assert.Equal("4", actual.field["wg"])
	assert.Equal("5", actual.field["wf"])
	assert.Equal("6", actual.field["uop_id"])
	assert.Equal("f", actual.field["stg"])
	assert.Equal("s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358",
		actual.field["asm"])

}

func TestParseInstructionExecute(t *testing.T) {
	assert := assert.New(t)
	parser := new(Parser)

	input := `si.inst id=1 cu=2 wf=3 uop_id=4 stg="su-r"`
	actual, err := parser.Parse(input)

	assert.Equal(nil, err)
	assert.Equal("si", actual.field["arch"])
	assert.Equal("1", actual.field["id"])
	assert.Equal("2", actual.field["cu"])
	assert.Equal("3", actual.field["wf"])
	assert.Equal("4", actual.field["uop_id"])
	assert.Equal("su-r", actual.field["stg"])
}

func TestParseInstructionEnd(t *testing.T) {
	assert := assert.New(t)
	parser := new(Parser)

	input := `si.end_inst id=1 cu=2`
	actual, err := parser.Parse(input)

	assert.Equal(nil, err)
	assert.Equal("si", actual.field["arch"])
	assert.Equal("1", actual.field["id"])
	assert.Equal("2", actual.field["cu"])
}

func TestParseMemoryNewAccess(t *testing.T) {
	assert := assert.New(t)
	parser := new(Parser)

	input := `mem.new_access name="A-1000" type="load" state="l1-cu02:load" addr=0xc610`
	actual, err := parser.Parse(input)

	assert.Equal(nil, err)
	assert.Equal("1000", actual.field["id"])
	assert.Equal("load", actual.field["type"])
	assert.Equal("l1-cu02", actual.field["module"])
	assert.Equal("load", actual.field["action"])
	assert.Equal("0xc610", actual.field["addr"])
}

func TestParseMemoryAccess(t *testing.T) {
	assert := assert.New(t)
	parser := new(Parser)

	input := `mem.access name="A-123" state="l1-cu0:find_and_lock"`
	actual, err := parser.Parse(input)

	assert.Equal(nil, err)
	assert.Equal("123", actual.field["id"])
	assert.Equal("l1-cu0", actual.field["module"])
	assert.Equal("find_and_lock", actual.field["action"])
}

func TestParseMemoryEndAccess(t *testing.T) {
	assert := assert.New(t)
	parser := new(Parser)

	input := `mem.end_access name="A-1000"`
	actual, err := parser.Parse(input)

	assert.Equal(nil, err)
	assert.Equal("1000", actual.field["id"])
}

func TestParseMemoryNewAccessBlock(t *testing.T) {
	assert := assert.New(t)
	parser := new(Parser)

	input := `mem.new_access_block cache="l2-4" access="A-64385" set=37 way=14`
	actual, err := parser.Parse(input)

	assert.Equal(nil, err)
	assert.Equal("l2-4", actual.field["cache"])
	assert.Equal("l2", actual.field["level"])
	assert.Equal("4", actual.field["module"])
	assert.Equal("A-64385", actual.field["id"])
	assert.Equal("37", actual.field["set"])
	assert.Equal("14", actual.field["way"])
}

func TestParseMemoryEndAccessBlock(t *testing.T) {
	assert := assert.New(t)
	parser := new(Parser)

	input := `mem.end_access_block cache="l2-3" access="A-16705" set=101 way=15`
	actual, err := parser.Parse(input)

	assert.Equal(nil, err)
	assert.Equal("l2-3", actual.field["cache"])
	assert.Equal("l2", actual.field["level"])
	assert.Equal("3", actual.field["module"])
	assert.Equal("A-16705", actual.field["id"])
	assert.Equal("101", actual.field["set"])
	assert.Equal("15", actual.field["way"])
}
