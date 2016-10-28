package trace

import "regexp"

// Parse the string and return a map of named group and the value
func parseNamedGroup(pattern, str string) map[string]string {
	r := regexp.MustCompile(pattern)
	match := r.FindStringSubmatch(str)
	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 {
			result[name] = match[i]
		}
	}
	return result
}

// ParseClock ...
// eg: c clk=1000
func ParseClock(raw string) map[string]string {
	return parseNamedGroup(`c clk=(?P<clock>\d+)`, raw)
}

// ParseInstructionUniqueID ...
// eg: id=69 cu=0
func ParseInstructionUniqueID(raw string) map[string]string {
	return parseNamedGroup(`id=(?P<id>\d+) cu=(?P<cu>\d+)`, raw)
}

// ParseInstructionNew ...
// eg: 'si.new_inst id=69 cu=0 ib=0 wg=0 wf=5 uop_id=8
// stg="f" asm="s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358"'
func ParseInstructionNew(raw string) map[string]string {
	return parseNamedGroup(`si.new_inst id=(?P<id>\d+) cu=(?P<cu>\d+) `+
		`ib=(?P<ib>\d+) wg=(?P<wg>\d+) wf=(?P<wf>\d+) uop_id=(?P<uop_id>\d+) `+
		`stg="(?P<stg>.+)" asm="(?P<asm>.*)"`, raw)
}

// ParseInstructionExecute ...
// eg: 'si.inst id=60 cu=0 wf=4 uop_id=7 stg="su-r"'
func ParseInstructionExecute(raw string) map[string]string {
	return parseNamedGroup(`si.inst id=(?P<id>\d+) cu=(?P<cu>\d+) `+
		`wf=(?P<wf>\d+) uop_id=(?P<uop_id>\d+) stg="(?P<stg>.+)"`, raw)
}

// ParseInstructionEnd ...
// eg: 'si.end_inst id=35 cu=3'
func ParseInstructionEnd(raw string) map[string]string {
	return parseNamedGroup(`si.end_inst id=(?P<id>\d+) cu=(?P<cu>\d+)`, raw)
}

// ParseMemoryUniqueID ...
// eg: 'name="A-227"'
func ParseMemoryUniqueID(raw string) map[string]string {
	return parseNamedGroup(`name="A-(?P<id>\d+)"`, raw)
}

// ParseMemoryNewAccess ...
// eg: mem.new_access name="A-227" type="load" state="l1-cu02:load" addr=0xc610
func ParseMemoryNewAccess(raw string) map[string]string {
	return parseNamedGroup(`mem.new_access name="A-(?P<id>\d+)" `+
		`type="(?P<type>\w+)" state="(?P<module>[^\"\:]+):(?P<action>[^\"\:]+)" `+
		`addr=(?P<addr>\w+)`, raw)
}

// ParseMemoryAccess ...
// eg: mem.access name="A-213" state="l1-cu0:find_and_lock"
func ParseMemoryAccess(raw string) map[string]string {
	return parseNamedGroup(`mem.access name="A-(?P<id>\d+)" `+
		`state="(?P<module>[^\"\:]+):(?P<action>\w+)"`, raw)
}

// ParseMemoryEndAccess ...
// eg: mem.end_access name="A-16512"
func ParseMemoryEndAccess(raw string) map[string]string {
	return parseNamedGroup(`mem.end_access name="A-(?P<id>\d+)"`, raw)
}

// ParseMemoryNewAccessBlock ...
// eg: mem.new_access_block cache="l2-4" access="A-64385" set=37 way=14
func ParseMemoryNewAccessBlock(raw string) map[string]string {
	return parseNamedGroup(`mem.new_access_block `+
		`cache="(?P<cache>(?P<level>\w+)-(?P<module>\w+))" `+
		`access="(?P<id>A-\d+)" set=(?P<set>\d+) way=(?P<way>\d+)`, raw)
}

// ParseMemoryEndAccessBlock ...
// eg: mem.end_access_block cache="l2-3" access="A-16705" set=101 way=15
func ParseMemoryEndAccessBlock(raw string) map[string]string {
	return parseNamedGroup(`mem.end_access_block `+
		`cache="(?P<cache>(?P<level>\w+)-(?P<module>\w+))" `+
		`access="(?P<id>A-\d+)" set=(?P<set>\d+) way=(?P<way>\d+)`, raw)
}
