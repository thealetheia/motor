package motor

import (
	"testing"

	"aletheia.icu/motor/assert"
)

func TestFmtexpr(t *testing.T) {
	assert := assert.New(t)

	expr := func(s string, tags ...Tag) tagmap {
		return tagmap{s, tags}
	}

	assert(expr("%d"), fmtexpr("%d"))
	assert(expr("%+v"), fmtexpr("%+v"))
	assert(expr("%+b", Tag{0, "label"}), fmtexpr("%{label}+b"))
	assert(expr("%.2f %2f", Tag{1, "label"}), fmtexpr("%.2f %{label}2f"))
}
