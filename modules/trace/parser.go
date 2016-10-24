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

// ParseClock example: 'c clk=1000'
func ParseClock(raw string) map[string]string {
	return parseNamedGroup(`c clk=(?P<clock>\d+)`, raw)
}

// ParseInstructionUniqueID example: 'id=69 cu=0'
func ParseInstructionUniqueID(raw string) map[string]string {
	return parseNamedGroup(`id=(?P<id>\d+) cu=(?P<cu>\d+)`, raw)
}

// ParseInstructionNew example: 'si.new_inst id=69 cu=0 ib=0 wg=0 wf=5 uop_id=8 stg="f" asm="s_load_dwordx4 s[8:11], s[2:3], 0x60 // 0000022C: C0880358"'
func ParseInstructionNew(raw string) map[string]string {
	return parseNamedGroup(`si.new_inst id=(?P<id>\d+) cu=(?P<cu>\d+) ib=(?P<ib>\d+) wg=(?P<wg>\d+) wf=(?P<wf>\d+) uop_id=(?P<uop_id>\d+) stg="(?P<stage>.+)" asm="(?P<asm>.*)"`, raw)
}
