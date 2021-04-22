package motor

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/sirupsen/logrus"
)

type mockWriter struct{}

func (w *mockWriter) Write(level int, args ...interface{}) {
	fmt.Fprintln(ioutil.Discard, args...)
}

func BenchmarkPrintClassic(b *testing.B) {
	f := ioutil.Discard
	for i := 0; i < b.N; i++ {
		fmt.Fprintln(f, "hello world")
		fmt.Fprintf(f, "interpolated %.2f string", 42.242)
		fmt.Fprint(f)
	}
}

func BenchmarkPrintMotor(b *testing.B) {
	log, _, _ := New(&mockWriter{})
	for i := 0; i < b.N; i++ {
		log("hello world")
		log("interpolated %.2f string", 42.242)
		log()
	}
}

func BenchmarkPrintLogrus(b *testing.B) {
	log := logrus.New()
	log.SetOutput(ioutil.Discard)
	for i := 0; i < b.N; i++ {
		log.Println("hello world")
		log.Printf("interpolated %.2f string", 42.242)
		log.Print()
	}
}
