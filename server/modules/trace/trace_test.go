package trace

import (
	"fmt"
	"testing"
)

func TestTrace(t *testing.T) {
	// TODO
	var trace Trace
	err := trace.Init("test/config.toml")
	if err != nil {
		fmt.Println(err)
	}
	trace.Process("test/test.gz")
}
