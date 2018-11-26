package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	// We want to be able to capture the output of our tracer in a bytes.Buffer so that wwe can then
	// ensure that the string in the buffer matches the expected value.
	var buf bytes.Buffer

	tracer := New(&buf)

	if tracer == nil {
		t.Error("Return from New should not be nil")
	} else {
		tracer.Trace("Hello trace package.")
		if buf.String() != "Hello trace package.\n" {
			t.Errorf("Trace should not write '%s'.", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer = Off()
	silentTracer.Trace("something")
}
