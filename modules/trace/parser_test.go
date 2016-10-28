package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseClock(t *testing.T) {
	assert := assert.New(t)
	testString := "c clk=1000"
	actual := ParseClock(testString)
	assert.Equal("1000", actual["clock"])
}

func TestParseInstructionUniqueID(t *testing.T) {
	assert := assert.New(t)
	testString := "id=1 cu=2"
	actual := ParseInstructionUniqueID(testString)
	assert.Equal("1", actual["id"])
	assert.Equal("2", actual["cu"])
}

func TestParseInstructionNew(t *testing.T) {
	assert := assert.New(t)
	testString := `si.new_inst id=1 cu=2 ib=3 wg=4 wf=5 uop_id=6 stg="f" ` +
		`asm="s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358"`
	actual := ParseInstructionNew(testString)
	assert.Equal("1", actual["id"])
	assert.Equal("2", actual["cu"])
	assert.Equal("3", actual["ib"])
	assert.Equal("4", actual["wg"])
	assert.Equal("5", actual["wf"])
	assert.Equal("6", actual["uop_id"])
	assert.Equal("f", actual["stg"])
	assert.Equal("s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358", actual["asm"])

}

func TestParseInstructionExecute(t *testing.T) {
	assert := assert.New(t)
	testString := `si.inst id=1 cu=2 wf=3 uop_id=4 stg="su-r"`
	actual := ParseInstructionExecute(testString)
	assert.Equal("1", actual["id"])
	assert.Equal("2", actual["cu"])
	assert.Equal("3", actual["wf"])
	assert.Equal("4", actual["uop_id"])
	assert.Equal("su-r", actual["stg"])
}

func TestParseInstructionEnd(t *testing.T) {
	assert := assert.New(t)
	testString := `si.end_inst id=1 cu=2`
	actual := ParseInstructionEnd(testString)
	assert.Equal("1", actual["id"])
	assert.Equal("2", actual["cu"])
}

func TestParseMemoryUniqueID(t *testing.T) {
	assert := assert.New(t)
	testString := `name="A-1000"`
	actual := ParseMemoryUniqueID(testString)
	assert.Equal("1000", actual["id"])
}

func TestParseMemoryNewAccess(t *testing.T) {
	assert := assert.New(t)
	testString := `mem.new_access name="A-1000" type="load" state="l1-cu02:load" addr=0xc610`
	actual := ParseMemoryNewAccess(testString)
	assert.Equal("1000", actual["id"])
	assert.Equal("load", actual["type"])
	assert.Equal("l1-cu02", actual["module"])
	assert.Equal("load", actual["action"])
	assert.Equal("0xc610", actual["addr"])
}

func TestParseMemoryAccess(t *testing.T) {
	assert := assert.New(t)
	testString := `mem.access name="A-123" state="l1-cu0:find_and_lock"`
	actual := ParseMemoryAccess(testString)
	assert.Equal("123", actual["id"])
	assert.Equal("l1-cu0", actual["module"])
	assert.Equal("find_and_lock", actual["action"])
}

func TestParseMemoryEndAccess(t *testing.T) {
	assert := assert.New(t)
	testString := `mem.end_access name="A-1000"`
	actual := ParseMemoryEndAccess(testString)
	assert.Equal("1000", actual["id"])
}

func TestParseMemoryNewAccessBlock(t *testing.T) {
	assert := assert.New(t)
	testString := `mem.new_access_block cache="l2-4" access="A-64385" set=37 way=14`
	actual := ParseMemoryNewAccessBlock(testString)
	assert.Equal("l2-4", actual["cache"])
	assert.Equal("l2", actual["level"])
	assert.Equal("4", actual["module"])
	assert.Equal("A-64385", actual["id"])
	assert.Equal("37", actual["set"])
	assert.Equal("14", actual["way"])
}

func TestParseMemoryEndAccessBlock(t *testing.T) {
	assert := assert.New(t)
	testString := `mem.end_access_block cache="l2-3" access="A-16705" set=101 way=15`
	actual := ParseMemoryEndAccessBlock(testString)
	assert.Equal("l2-3", actual["cache"])
	assert.Equal("l2", actual["level"])
	assert.Equal("3", actual["module"])
	assert.Equal("A-16705", actual["id"])
	assert.Equal("101", actual["set"])
	assert.Equal("15", actual["way"])
}
