package trace

import "testing"

func TestParseClock(t *testing.T) {
	testString := "c clk=1000"
	result := ParseClock(testString)
	if result["clock"] == "1000" {
		return
	}
	t.Error("ParseClock failed")
}

func TestParseInstructionUniqueID(t *testing.T) {
	testString := "id=69 cu=0"
	result := ParseInstructionUniqueID(testString)
	if result["id"] == "69" && result["cu"] == "0" {
		return
	}
	t.Error("ParseInstructionUniqueID failed")
}

func TestParseInstructionNew(t *testing.T) {
	testString := `si.new_inst id=1 cu=2 ib=3 wg=4 wf=5 uop_id=6 stg="f" asm="s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358"`
	result := ParseInstructionNew(testString)
	if result["id"] != "1" {
		t.Errorf("RETURN = %s, EXPECT 1", result["id"])
	}
	if result["cu"] != "2" {
		t.Errorf("RETURN = %s, EXPECT 2", result["cu"])
	}
	if result["ib"] != "3" {
		t.Errorf("RETURN = %s, EXPECT 3", result["ib"])
	}
	if result["wg"] != "4" {
		t.Errorf("RETURN = %s, EXPECT 4", result["wg"])
	}
	if result["wf"] != "5" {
		t.Errorf("RETURN = %s, EXPECT 5", result["wf"])
	}
	if result["uop_id"] != "6" {
		t.Errorf("RETURN = %s, EXPECT 6", result["uop_id"])
	}
	if result["stage"] != "f" {
		t.Errorf("RETURN = %s, EXPECT f", result["stage"])
	}
	if result["asm"] != "s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358" {
		t.Errorf("RETURN = %s, EXPECT sload_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358|", result["asm"])
	}
	return
}
