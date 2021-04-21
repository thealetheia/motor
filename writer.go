package motor

import (
	"fmt"
	"os"
)

type Writer interface {
	Write(level int, args ...interface{})
}

type stdoutWriter struct{}

func (w *stdoutWriter) Write(level int, args ...interface{}) {
	fmt.Fprintln(os.Stdout, args...)
}

type stderrWriter struct{}

func (w *stderrWriter) Write(level int, args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
}
