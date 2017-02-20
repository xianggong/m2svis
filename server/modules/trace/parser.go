package trace

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// Parser for trace
type Parser struct {
	currCycle int
}

// ParseInfo contains cycle, key and field information of each line
type ParseInfo struct {
	cycle int
	key   string
	field map[string]string
}

// ParserRegexTable is a table containing string to pattern string
// For input string `si.inst id=1 cu=2 wf=3 uop_id=4 stg="su-r"`
// ParserRegexTable["si.inst"] returns the regex pattern
var parserRegexTable = map[string]string{
	`c`:                    `c clk=(?P<clock>\d+)`,
	`si.new_inst`:          `(?P<arch>\w+).new_inst id=(?P<id>\d+) cu=(?P<cu>\d+) ib=(?P<ib>\d+) wg=(?P<wg>\d+) wf=(?P<wf>\d+) uop_id=(?P<uop_id>\d+) stg="(?P<stg>.+)" asm="(?P<asm>.*)"`,
	`si.inst`:              `(?P<arch>\w+).inst id=(?P<id>\d+) cu=(?P<cu>\d+) wf=(?P<wf>\d+) uop_id=(?P<uop_id>\d+) stg="(?P<stg>.+)"`,
	`si.end_inst`:          `(?P<arch>\w+).end_inst id=(?P<id>\d+) cu=(?P<cu>\d+)`,
	`mem.new_access`:       `mem.new_access name="A-(?P<id>\d+)" type="(?P<type>\w+)" state="(?P<module>[^\"\:]+):(?P<action>[^\"\:]+)" addr=(?P<addr>\w+)`,
	`mem.access`:           `mem.access name="A-(?P<id>\d+)" state="(?P<module>[^\"\:]+):(?P<action>\w+)"`,
	`mem.end_access`:       `mem.end_access name="A-(?P<id>\d+)"`,
	`mem.new_access_block`: `mem.new_access_block cache="(?P<cache>(?P<level>\w+)-(?P<module>\w+))" access="(?P<id>A-\d+)" set=(?P<set>\d+) way=(?P<way>\d+)`,
	`mem.end_access_block`: `mem.end_access_block cache="(?P<cache>(?P<level>\w+)-(?P<module>\w+))" access="(?P<id>A-\d+)" set=(?P<set>\d+) way=(?P<way>\d+)`,
}

// Parse the string and return a map of named group and the value
func parseNamedGroup(pattern, str string) map[string]string {
	r := regexp.MustCompile(pattern)
	match := r.FindStringSubmatch(str)

	// If find match, return each field and value as a map
	if match != nil {
		result := make(map[string]string)
		for i, name := range r.SubexpNames() {
			if i != 0 {
				result[name] = match[i]
			}
		}
		return result
	}

	// Otherwise return nil
	return nil
}

// Parse takes an input string and return parse result and error
func (parser *Parser) Parse(input string) (pr ParseInfo, err error) {
	// Get identifier
	key := strings.Fields(input)[0]

	// Get pattern in the regex table
	pattern := parserRegexTable[key]
	if pattern == "" {
		return ParseInfo{}, errors.New("No Pattern")
	}

	// Get parse result
	result := parseNamedGroup(pattern, input)

	// Update current cycle
	if key == `c` {
		parser.currCycle, _ = strconv.Atoi(result["clock"])
	}

	return ParseInfo{cycle: parser.currCycle, key: key, field: result}, nil
}

// GetID returns id for memory info and id_cu for instruction info
func (parseInfo *ParseInfo) GetID() string {
	switch parseInfo.key {
	case "si.new_inst", "si.inst", "si.end_inst":
		return parseInfo.field["id"] + "_" + parseInfo.field["cu"]
	case "mem.new_access", "mem.access", "mem.end_access",
		"mem.new_access_block", "mem.end_access_block":
		return parseInfo.field["id"]
	}
	return ""
}
