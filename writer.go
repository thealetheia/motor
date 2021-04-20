package motor

import "fmt"

type Writer interface {
	Write(level int, args ...interface{})
}

type stdWriter struct{}

func (w *stdWriter) Write(level int, args ...interface{}) {
	fmt.Println(args...)
}
