package goolog2

import (
	"io"
)

// Holder of a file
//
// This interface wraps a logging file. The file can be just one physical
// file. Or it can be a pattern defining a rotating logs.
type FileHolder interface {
	// Access and log the log writer for writing of one item
	AccessWriter(
		functor func(writer io.Writer))

	// Close the file holder
	Close()
}
