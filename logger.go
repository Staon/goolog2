package goolog2

import ()

// An object logger
type Logger interface {
	// Destroy the logger
	//
	// Warning: the method is not thread safe!
	Destroy()

	// Log an object
	//
	// Parameters:
	//     system: a logging system
	//     subsystem: a logging sybsystem
	//     severity: message severity
	//     verbosity: message verbosity
	//     object: an object to be logged
	LogObject(
		system string,
		subsystem Subsystem,
		severity Severity,
		verbosity Verbosity,
		object interface{})
}
