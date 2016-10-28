package trace

import "testing"
import "github.com/stretchr/testify/assert"

var inst = Instruction{}

func TestInstructionNew(t *testing.T) {
	assert := assert.New(t)
	testString := `si.new_inst id=1 cu=2 ib=3 wg=4 wf=5 uop_id=6 stg="f" ` +
		`asm="s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358"`
	actual := ParseInstructionNew(testString)
	var cycle = 1000
	inst.New(cycle, actual)
	assert.Equal(1, inst.id)
	assert.Equal(1000, inst.start)
	assert.Equal(0, inst.finish)
	assert.Equal(0, inst.length)
	assert.Equal(0, inst.fetch)
	assert.Equal(0, inst.decode)
	assert.Equal(0, inst.issue)
	assert.Equal(0, inst.execute)
	assert.Equal(0, inst.write)
	assert.Equal(2, inst.cu)
	assert.Equal(3, inst.ib)
	assert.Equal(4, inst.wg)
	assert.Equal(5, inst.wf)
	assert.Equal(6, inst.uop)
	assert.Equal("s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358", inst.asm)
	assert.Equal(Activity{1000, "f"}, inst.lifeVerbose[0])
	assert.Empty(inst.lifeConcise)
}

func TestInstructionExe(t *testing.T) {
	TestInstructionNew(t)

	assert := assert.New(t)
	testString := `si.inst id=1 cu=2 wf=5 uop_id=6 stg="su-r"`
	actual := ParseInstructionExecute(testString)
	cycle := 2000
	inst.Exe(cycle, actual)

	assert.Equal(1, inst.id)
	assert.Equal(1000, inst.start)
	assert.Equal(0, inst.finish)
	assert.Equal(0, inst.length)
	assert.Equal(0, inst.fetch)
	assert.Equal(0, inst.decode)
	assert.Equal(0, inst.issue)
	assert.Equal(0, inst.execute)
	assert.Equal(0, inst.write)
	assert.Equal(2, inst.cu)
	assert.Equal(3, inst.ib)
	assert.Equal(4, inst.wg)
	assert.Equal(5, inst.wf)
	assert.Equal(6, inst.uop)
	assert.Equal("s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358", inst.asm)
	assert.Equal(Activity{2000, "su-r"}, inst.lifeVerbose[2])
	assert.Equal(Activity{1000, "f"}, inst.lifeConcise[0])
}
