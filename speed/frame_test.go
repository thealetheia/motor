package speed

import (
	"fmt"
	"testing"
	"time"
)

var log = fmt.Println

func TestFrameFmt(t *testing.T) {
	var (
		e = time.Now()
		b = time.Now().Add(-50 * time.Millisecond)
	)

	f := F{N: 5, Begin: b, End: e}
	fmt.Printf("%.95f\n", f)
}
