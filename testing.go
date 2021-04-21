package motor

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/fatih/color"
)

func init() {
	const debug = "!probe.debug"
	const trace = "!probe.trace"

	k := -1
	for i, arg := range os.Args {
		if arg == debug || arg == trace {
			k = i
			break
		}
	}
	if k < 0 {
		return
	}
	Args = os.Args[k+1:]
	switch os.Args[k] {
	case debug:
		V(2)
	case trace:
		V(3)
	}
}

// Assertf functions as
func Assertf(t *testing.T) func(expected interface{}, format string, values ...interface{}) {
	return func(l interface{}, f string, r ...interface{}) {
		left := fmt.Sprint(l)
		right := fmt.Sprintf(f, r...)
		if left == right {
			return
		}
		_, file, line, _ := runtime.Caller(1)

		bold := color.New(color.Bold)
		redBold := color.New(color.FgRed, color.Bold)
		redBold.Print("!!! ")
		fmt.Printf("%v:%v\n", file, line)
		if diff := diff(left, right); diff != "" {
			fmt.Printf(diff)
		} else {
			bold.Println("\tEXPECTED")
			fmt.Println(left)
			bold.Println("\tRECEIVED")
			fmt.Println(right, "\a")
		}
		t.FailNow()
	}
}
