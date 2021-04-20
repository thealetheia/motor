package speed

import "errors"

func ø(err string) error {
	return errors.New("speed: " + err)
}

var (
	ErrFrameClosed = ø("calling end to a closed frame")
	ErrLateFrame   = ø("fulfilled frame is late")
)
