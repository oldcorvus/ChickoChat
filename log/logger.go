package logger

import (
	"fmt"
	"io"
)

// logging events throughout code.
type Logger interface {
	Log(...interface{})
}

// New creates a new Logger that will write the output to
// the specified io.Writer.
func New(w io.Writer) Logger {
	return &logger{out: w}
}

type logger struct {
	out io.Writer
}

// writes the arguments io.Writer.
func (t *logger) Log(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}
