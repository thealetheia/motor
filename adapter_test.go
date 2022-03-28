package motor

import (
	"html"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"aletheia.icu/motor/speed"
)

func randomEmoji() string {
	rand.Seed(time.Now().UnixNano())
	// http://apps.timwhitlock.info/emoji/tables/unicode
	emoji := [][]int{
		// Emoticons icons
		{128513, 128591},
		// Transport and map symbols
		{128640, 128704},
	}
	r := emoji[rand.Int()%len(emoji)]
	min := r[0]
	max := r[1]
	n := rand.Intn(max-min+1) + min
	return html.UnescapeString("&#" + strconv.Itoa(n) + ";")
}

func TestSimple(t *testing.T) {
	// a simple motor
	m := New(Config{Sinks: []Adapter{&Simple{}}, Debug: true})

	brr := m.Func("/spool", randomEmoji())
	defer brr.Flush()

	for i := 0; i < 5; i++ {
		<-speed.After(100*i, "us")
		if i%2 == 0 {
			brr.Printf("hello %(world)d", i+1)
		} else {
			brr.Debugf("hello %(world)d", i+1)
		}
	}
}
