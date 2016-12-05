package trace

import "testing"

func TestTrace(t *testing.T) {
	// TODO
	var trace Trace
	trace.Init("test/config.toml")
	trace.Process("test/test.gz")
}
