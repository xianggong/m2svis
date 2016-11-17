package trace

import (
	"bufio"
	"compress/gzip"
	"log"
	"os"
)

// Trace contains instruction pool and parser object
type Trace struct {
	instPool InstructionPool
	parser   Parser
}

// Process trace
func (trace *Trace) Process(path string) {
	// Get trace file
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	// Open as gzip
	gzfile, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer gzfile.Close()

	// New scanner to read file line by line
	scanner := bufio.NewScanner(gzfile)
	scanner.Split(bufio.ScanLines)

	// Read line by line and parse
	for scanner.Scan() {
		line := scanner.Text()

		// check for errors
		if err = scanner.Err(); err != nil {
			log.Fatal(err)
		}

		info, err := trace.parser.Parse(line)
		if err == nil {
			trace.instPool.Receive(&info)
		}
	}
}

// GetJSON returns JSON arrays
func (trace *Trace) GetJSON() []*InstructionJSON {
	if len(trace.instPool.Ready) == 0 {
		return nil
	}

	var traceJSON []*InstructionJSON

	for _, inst := range trace.instPool.Ready {
		traceJSON = append(traceJSON, inst.GetOverviewJSON()...)
	}

	return traceJSON
}
