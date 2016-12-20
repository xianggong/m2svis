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
	assert.Equal(1, inst.ID)
	assert.Equal(1000, inst.Start)
	assert.Equal(0, inst.Finish)
	assert.Equal(0, inst.Length)
	assert.Equal(2, inst.CU)
	assert.Equal(3, inst.IB)
	assert.Equal(4, inst.WG)
	assert.Equal(5, inst.WF)
	assert.Equal(6, inst.UOP)
	assert.Equal("[4-5]: s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358", inst.Assembly)
	assert.Equal(Activity{1000, "f"}, inst.LifeVerbose[0])
	assert.Equal(Activity{0, "f"}, inst.LifeConcise[0])
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
	assert.Equal(1, inst.ID)
	assert.Equal(1000, inst.Start)
	assert.Equal(0, inst.Finish)
	assert.Equal(0, inst.Length)
	assert.Equal(2, inst.CU)
	assert.Equal(3, inst.IB)
	assert.Equal(4, inst.WG)
	assert.Equal(5, inst.WF)
	assert.Equal(6, inst.UOP)
	assert.Equal("[4-5]: s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358", inst.Assembly)
	assert.Equal(Activity{1000, "f"}, inst.LifeVerbose[0])
	assert.Equal(Activity{1500, "su-r"}, inst.LifeVerbose[1])
	assert.Equal(Activity{500, "f"}, inst.LifeConcise[0])
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
	assert.Equal(1, inst.ID)
	assert.Equal(1000, inst.Start)
	assert.Equal(0, inst.Finish)
	assert.Equal(0, inst.Length)
	assert.Equal(2, inst.CU)
	assert.Equal(3, inst.IB)
	assert.Equal(4, inst.WG)
	assert.Equal(5, inst.WF)
	assert.Equal(6, inst.UOP)
	assert.Equal("[4-5]: s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358", inst.Assembly)
	assert.Equal(Activity{1000, "f"}, inst.LifeVerbose[0])
	assert.Equal(Activity{1500, "su-r"}, inst.LifeVerbose[1])
	assert.Equal(Activity{500, "f"}, inst.LifeConcise[0])

	// Input
	inputEnd := map[string]string{}
	inputEnd["id"] = "1"
	inputEnd["cu"] = "2"

	// End
	inst.End(2000, inputEnd)

	// Test body
	assert.Equal(1000, inst.Start)
	assert.Equal(2000, inst.Finish)
	assert.Equal(1000, inst.Length)
	assert.Equal(3, len(inst.LifeVerbose))
	assert.Equal(Activity{1000, "f"}, inst.LifeVerbose[0])
	assert.Equal(Activity{1500, "su-r"}, inst.LifeVerbose[1])
	assert.Equal(Activity{2000, "end"}, inst.LifeVerbose[2])
	assert.Equal(2, len(inst.LifeConcise))
	assert.Equal(Activity{500, "f"}, inst.LifeConcise[0])
	assert.Equal(Activity{500, "su-r"}, inst.LifeConcise[1])
}
