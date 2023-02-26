package logger

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	logger := New(&buf)
	if logger == nil {
		t.Error("should not be nil")
	} else {
		logger.Log("log package.")
		if buf.String() != "log package.\n" {
			t.Errorf("Trace should not write '%s'.", buf.String())
		}
	}
}
