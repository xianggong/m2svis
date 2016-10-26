package trace

import "testing"

func testField(t *testing.T, result map[string]string, name, expect string) {
	if result[name] != expect {
		t.Errorf("RETURN = %s, EXPECT %s", result[name], expect)
	}
}

func TestParseClock(t *testing.T) {
	testString := "c clk=1000"
	result := ParseClock(testString)
	testField(t, result, "clock", "1000")
}

func TestParseInstructionUniqueID(t *testing.T) {
	testString := "id=1 cu=2"
	result := ParseInstructionUniqueID(testString)
	testField(t, result, "id", "1")
	testField(t, result, "cu", "2")
}

func TestParseInstructionNew(t *testing.T) {
	testString := `si.new_inst id=1 cu=2 ib=3 wg=4 wf=5 uop_id=6 stg="f" asm="s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358"`
	result := ParseInstructionNew(testString)
	testField(t, result, "id", "1")
	testField(t, result, "cu", "2")
	testField(t, result, "ib", "3")
	testField(t, result, "wg", "4")
	testField(t, result, "wf", "5")
	testField(t, result, "uop_id", "6")
	testField(t, result, "stg", "f")
	testField(t, result, "asm", "s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358")

}

func TestParseInstructionExecute(t *testing.T) {
	testString := `si.inst id=1 cu=2 wf=3 uop_id=4 stg="su-r"`
	result := ParseInstructionExecute(testString)
	testField(t, result, "id", "1")
	testField(t, result, "cu", "2")
	testField(t, result, "wf", "3")
	testField(t, result, "uop_id", "4")
	testField(t, result, "stg", "su-r")
}

func TestParseInstructionEnd(t *testing.T) {
	testString := `si.end_inst id=1 cu=2`
	result := ParseInstructionEnd(testString)
	testField(t, result, "id", "1")
	testField(t, result, "cu", "2")
}

func TestParseMemoryUniqueID(t *testing.T) {
	testString := `name="A-1000"`
	result := ParseMemoryUniqueID(testString)
	testField(t, result, "id", "1000")
}

func TestParseMemoryNewAccess(t *testing.T) {
	testString := `mem.new_access name="A-1000" type="load" state="l1-cu02:load" addr=0xc610`
	result := ParseMemoryNewAccess(testString)
	testField(t, result, "id", "1000")
	testField(t, result, "type", "load")
	testField(t, result, "module", "l1-cu02")
	testField(t, result, "action", "load")
	testField(t, result, "addr", "0xc610")
}

func TestParseMemoryAccess(t *testing.T) {
	testString := `mem.access name="A-123" state="l1-cu0:find_and_lock"`
	result := ParseMemoryAccess(testString)
	testField(t, result, "id", "123")
	testField(t, result, "module", "l1-cu0")
	testField(t, result, "action", "find_and_lock")
}

func TestParseMemoryEndAccess(t *testing.T) {
	testString := `mem.end_access name="A-1000"`
	result := ParseMemoryEndAccess(testString)
	testField(t, result, "id", "1000")
}

func TestParseMemoryNewAccessBlock(t *testing.T) {
	testString := `mem.new_access_block cache="l2-4" access="A-64385" set=37 way=14`
	result := ParseMemoryNewAccessBlock(testString)
	testField(t, result, "cache", "l2-4")
	testField(t, result, "level", "l2")
	testField(t, result, "module", "4")
	testField(t, result, "id", "A-64385")
	testField(t, result, "set", "37")
	testField(t, result, "way", "14")
}

func TestParseMemoryEndAccessBlock(t *testing.T) {
	testString := `mem.end_access_block cache="l2-3" access="A-16705" set=101 way=15`
	result := ParseMemoryEndAccessBlock(testString)
	testField(t, result, "cache", "l2-3")
	testField(t, result, "level", "l2")
	testField(t, result, "module", "3")
	testField(t, result, "id", "A-16705")
	testField(t, result, "set", "101")
	testField(t, result, "way", "15")
}
