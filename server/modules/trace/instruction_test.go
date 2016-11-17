package trace

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// clk=1000 si.new_inst id=1 cu=2 ib=3 wg=4 wf=5 uop_id=6 stg="f" asm="s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358"
func TestInstructionNew(t *testing.T) {
	assert := assert.New(t)
	inst := new(Instruction)

	// Input
	inputNew := map[string]string{}
	inputNew["id"] = "1"
	inputNew["cu"] = "2"
	inputNew["ib"] = "3"
	inputNew["wg"] = "4"
	inputNew["wf"] = "5"
	inputNew["uop_id"] = "6"
	inputNew["stg"] = "f"
	inputNew["asm"] = "s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358"

	// New
	inst.New(1000, inputNew)

	// Test body
	assert.Equal(1, inst.id)
	assert.Equal(1000, inst.start)
	assert.Equal(0, inst.finish)
	assert.Equal(0, inst.length)
	assert.Equal(2, inst.cu)
	assert.Equal(3, inst.ib)
	assert.Equal(4, inst.wg)
	assert.Equal(5, inst.wf)
	assert.Equal(6, inst.uop)
	assert.Equal("s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358", inst.asm)
	assert.Equal(Activity{1000, "f"}, inst.lifeVerbose[0])
	assert.Equal(Activity{0, "f"}, inst.lifeConcise[0])
}

// clk=1000 si.new_inst id=1 cu=2 ib=3 wg=4 wf=5 uop_id=6 stg="f" asm="s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358"
// clk=1500 si.inst id=1 cu=2 wf=5 uop_id=6 stg="su-r"
func TestInstructionExe(t *testing.T) {
	assert := assert.New(t)
	inst := new(Instruction)

	// Input
	inputNew := map[string]string{}
	inputNew["id"] = "1"
	inputNew["cu"] = "2"
	inputNew["ib"] = "3"
	inputNew["wg"] = "4"
	inputNew["wf"] = "5"
	inputNew["uop_id"] = "6"
	inputNew["stg"] = "f"
	inputNew["asm"] = "s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358"

	// New
	inst.New(1000, inputNew)

	// Input
	inputExe := map[string]string{}
	inputExe["id"] = "1"
	inputExe["cu"] = "2"
	inputExe["wf"] = "5"
	inputExe["uop_id"] = "6"
	inputExe["stg"] = "su-r"

	// Execute
	inst.Exe(1500, inputExe)

	// Test body
	assert.Equal(1, inst.id)
	assert.Equal(1000, inst.start)
	assert.Equal(0, inst.finish)
	assert.Equal(0, inst.length)
	assert.Equal(2, inst.cu)
	assert.Equal(3, inst.ib)
	assert.Equal(4, inst.wg)
	assert.Equal(5, inst.wf)
	assert.Equal(6, inst.uop)
	assert.Equal("s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358", inst.asm)
	assert.Equal(Activity{1000, "f"}, inst.lifeVerbose[0])
	assert.Equal(Activity{1500, "su-r"}, inst.lifeVerbose[1])
	assert.Equal(Activity{500, "f"}, inst.lifeConcise[0])
}

// clk=1000 si.new_inst id=1 cu=2 ib=3 wg=4 wf=5 uop_id=6 stg="f" asm="s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358"
// clk=1500 si.inst id=1 cu=2 wf=5 uop_id=6 stg="su-r"
// clk=2000 si.end id=1 cu=2
func TestInstructionEnd(t *testing.T) {
	assert := assert.New(t)
	inst := new(Instruction)

	// Input
	inputNew := map[string]string{}
	inputNew["id"] = "1"
	inputNew["cu"] = "2"
	inputNew["ib"] = "3"
	inputNew["wg"] = "4"
	inputNew["wf"] = "5"
	inputNew["uop_id"] = "6"
	inputNew["stg"] = "f"
	inputNew["asm"] = "s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358"

	// New
	inst.New(1000, inputNew)

	// Input
	inputExe := map[string]string{}
	inputExe["id"] = "1"
	inputExe["cu"] = "2"
	inputExe["wf"] = "5"
	inputExe["uop_id"] = "6"
	inputExe["stg"] = "su-r"

	// Execute
	inst.Exe(1500, inputExe)

	// Test body
	assert.Equal(1, inst.id)
	assert.Equal(1000, inst.start)
	assert.Equal(0, inst.finish)
	assert.Equal(0, inst.length)
	assert.Equal(2, inst.cu)
	assert.Equal(3, inst.ib)
	assert.Equal(4, inst.wg)
	assert.Equal(5, inst.wf)
	assert.Equal(6, inst.uop)
	assert.Equal("s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358", inst.asm)
	assert.Equal(Activity{1000, "f"}, inst.lifeVerbose[0])
	assert.Equal(Activity{1500, "su-r"}, inst.lifeVerbose[1])
	assert.Equal(Activity{500, "f"}, inst.lifeConcise[0])

	// Input
	inputEnd := map[string]string{}
	inputEnd["id"] = "1"
	inputEnd["cu"] = "2"

	// End
	inst.End(2000, inputEnd)

	// Test body
	assert.Equal(1000, inst.start)
	assert.Equal(2000, inst.finish)
	assert.Equal(1000, inst.length)
	assert.Equal(3, len(inst.lifeVerbose))
	assert.Equal(Activity{1000, "f"}, inst.lifeVerbose[0])
	assert.Equal(Activity{1500, "su-r"}, inst.lifeVerbose[1])
	assert.Equal(Activity{2000, "end"}, inst.lifeVerbose[2])
	assert.Equal(2, len(inst.lifeConcise))
	assert.Equal(Activity{500, "f"}, inst.lifeConcise[0])
	assert.Equal(Activity{500, "su-r"}, inst.lifeConcise[1])
}

func TestInstructionGetJSON(t *testing.T) {
	// TODO
}
