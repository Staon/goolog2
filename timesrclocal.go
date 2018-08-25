package goolog2

import (
	"time"
)

type timeSourceLocal struct {
}

// Create new local time source
//
// This time source returns local current time read from the system
func NewTimeSourceLocal() TimeSource {
	return &timeSourceLocal{}
}

func (this *timeSourceLocal) Now() time.Time {
	return time.Now()
}
