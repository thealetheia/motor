package motor

import (
	"testing"

	"aletheia.icu/motor/assert"
)

func TestFmtexpr(t *testing.T) {
	assert := assert.New(t)

	assert(tagmap{"%d", nil}, fmtexpr("%d"))
	assert(tagmap{"%+v", nil}, fmtexpr("%+v"))
	assert(tagmap{"%+b", []tag{
		{0, "label"},
	}}, fmtexpr("%{label}+b"))
	assert(tagmap{"%.2f %2f", []tag{
		{1, "label"},
	}}, fmtexpr("%.2f %{label}2f"))
}
