package trace

import (
	"fmt"
	"io"
)

// Tracer is the interface that describes an object capable of
// tracing events througout code
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

// Accepting io.Writer means that the user can decide where the tracing output
// will be writter. This output could be the standard output, a file, network socket,
// bytes.Buffer as in our test case.
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

// Notice nilTracer does not take io.Writer like the tracer struct above.
// It does not need one because it doesn't need to write anything.
type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// Off creates a Tracer that will ignore calls to Trace.
func Off() Tracer {
	return &nilTracer{}
}
