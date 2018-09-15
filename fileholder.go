package goolog2

import (
	"io"
)

// Color constants
type Color int

const (
	NONE Color = iota
	RED
	YELLOW
	BLUE
)

// Generic file writer
//
// This interface is a classic io.Writer interface extended
// by possibility to change colors of the output.
type FileWriter interface {
	io.Writer

	// Flush the output
	Sync()

	// Change output color
	ChangeColor(
		color Color)

	// Reset the output color
	ResetColor()
}

// Holder of a file
//
// This interface wraps a logging file. The file can be just one physical
// file. Or it can be a pattern defining a rotating logs.
type FileHolder interface {
	// Access and log the log writer for writing of one item
	AccessWriter(
		functor func(writer FileWriter))

	// Close the file holder
	Close()
}
