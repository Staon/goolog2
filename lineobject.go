package goolog2

import ()

// Simple line logging object
//
// This interface wraps simple logging message - one line of text
type LineObject interface {
	// Get the message line
	GetLogLine() string
}
