package assert

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/google/go-cmp/cmp"
)

type AssertFn = func(assertion ...interface{})

var (
	bold    = color.New(color.Bold)
	redBold = color.New(color.FgRed, color.Bold)
)

func For(t *testing.T) AssertFn {
	return func(A ...interface{}) {
		var b bytes.Buffer
		switch len(A) {
		// silent
		case 0:
		// error, bool
		case 1:
			switch a := A[0].(type) {
			case error:
				if a == nil {
					return
				}
			case bool:
				if a == true {
					return
				}
			default:
				if a != nil {
					return
				}
			}
		// pair-wise
		case 2:
			if cmp.Equal(A[0], A[1]) {
				return
			}
			b.WriteString(cmp.Diff(A[0], A[1]))
		// fmt style
		default:
			left := fmt.Sprint(A[0])
			format := A[1].(string)
			right := fmt.Sprintf(format, A[2:]...)
			if left == right {
				return
			}
			if diff := diff(left, right); diff != "" {
				fmt.Fprintf(&b, diff)
			} else {
				bold.Fprintln(&b, "\tEXPECTED")
				fmt.Fprintln(&b, left)
				bold.Fprintln(&b, "\tRECEIVED")
				fmt.Fprintln(&b, right, "\a")
			}
		}

		localise()
		if b.Len() > 0 {
			fmt.Print(b.String())
		}
		// runtime.Breakpoint()
		t.FailNow()
	}
}

func localise() {
	prepare := func(s string) string {
		return "\t" + strings.Trim(s, "\t ")
	}

	_, file, line, _ := runtime.Caller(2)
	line--
	fb, _ := ioutil.ReadFile(file)
	lines := strings.Split(string(fb), "\n")

	j := 0
	for i := line - 1; i > 0; i-- {
		lines[i] = prepare(lines[i])
		if !strings.HasPrefix(lines[i], "\t//") {
			break
		}
		j++
	}

	redBold.Print("!!! ")
	fmt.Printf("%v:%v\n", file, line+1)

	fmt.Println()
	if j > 0 {
		comment := strings.Join(lines[line-j:line], "\n")
		fmt.Println(comment)
	}
	fmt.Println(prepare(lines[line]))
	fmt.Println()
}
