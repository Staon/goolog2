package goolog2

import (
	"io"
	"os"
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
	// If writter is not closable, the method Close returns nil.
	io.WriteCloser

	// If it is unable to read file info, it returns nil
	Stat() os.FileInfo

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

	// Increase reference counter - there is new owner of the holder
	//
	// Returns: itself
	Ref() FileHolder

	// Decrease reference counter
	//
	// If the counter reaches zero the object is destroyed (files
	// closed etc.)
	//
	// Expectation: there are no other threads accessing the holder
	//     if the reference counter reaches zero!
	Unref()
}
