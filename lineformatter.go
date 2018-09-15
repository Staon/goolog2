package goolog2

import (
	"time"
)

// Generic line formatter
//
// Implementations of this class format output logging lines
type LineFormatter interface {
	// Format a line
	//
	// Parameters:
	//     writer: line writer
	//     now: current time
	//     system: logging system
	//     subsystem: logging subsystem (can be empty)
	//     severity: severity of the logging message
	//     verbosity: verbosity of the logging message
	//     line: the logging message
	FormatLine(
		writer FileWriter,
		now time.Time,
		system string,
		subsystem Subsystem,
		severity Severity,
		verbosity Verbosity,
		line string)
}
