package goolog2

import (
	"time"
)

// Source of current time used in the logging messages
type TimeSource interface {
	// Get current time
	Now() time.Time
}
