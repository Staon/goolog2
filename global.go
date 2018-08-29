package goolog2

import (
	"os"
)

var globalLog LogDispatcher
var timeSource TimeSource

// Initialize the global log
//
// Parameters:
//     system: a system identifier shown in the logs
func Init(
	system string,
) {
	timeSource = NewTimeSourceLocal()
	globalLog = NewLogDispatcher(system)
}

// Initialize the global log with specified time source
func InitWithTimeSource(
	system string,
	timesrc TimeSource,
) {
	timeSource = timesrc
	globalLog = NewLogDispatcher(system)
}

// Destroy the global log
func Destroy() {
	globalLog.Destroy()
	globalLog = nil
	timeSource = nil
}

// Add a logger
//
// Parameters:
//     name: ID of the logger
//     subsystem: logging subsystem. Can be empty.
//     severities: mask of logging severities
//     verbosity: logging verbosity
//     logger: the logger. The ownership is taken - the logger
//             is destroyed with the dispatcher.
func AddLogger(
	name string,
	subsystem Subsystem,
	severities SeverityMask,
	verbosity Verbosity,
	logger Logger,
) {
	globalLog.AddLogger(name, subsystem, severities, verbosity, logger)
}

// Add a simple file logger
//
// Parameters:
//     name: ID of the logger
//     subsystem: logging subsystem. Can be empty.
//     severities: mask of logging severities
//     verbosity: logging verbosity
//     file: path to the logging file
//     sync: flush all message immediately
func AddFileLogger(
	name string,
	subsystem Subsystem,
	severities SeverityMask,
	verbosity Verbosity,
	file string,
	sync bool,
) {
	f := NewSimpleFile(file, sync)
	logger := NewFileLogger(timeSource, f)
	AddLogger(name, subsystem, severities, verbosity, logger)
}

// Add a pattern file logger
//
// Parameters:
//     name: ID of the logger
//     subsystem: logging subsystem. Can be empty.
//     severities: mask of logging severities
//     verbosity: logging verbosity
//     pattern: pattern of names the log files
//     sync: flush all message immediately
func AddPatternFileLogger(
	name string,
	subsystem Subsystem,
	severities SeverityMask,
	verbosity Verbosity,
	pattern string,
	sync bool,
) {
	f := NewPatternFile(timeSource, pattern, sync)
	logger := NewFileLogger(timeSource, f)
	AddLogger(name, subsystem, severities, verbosity, logger)
}

// Add a console logger
//
// Parameters:
//     name: ID of the logger
//     subsystem: logging subsystem. Can be empty.
//     severities: mask of logging severities
//     verbosity: logging verbosity
//     output: an output stream
func AddConsoleLogger(
	name string,
	subsystem Subsystem,
	severities SeverityMask,
	verbosity Verbosity,
	output *os.File,
) {
	logger := NewConsoleLogger(output)
	AddLogger(name, subsystem, severities, verbosity, logger)
}

// Add a console logger on the standard error
//
// Parameters:
//     name: ID of the logger
//     subsystem: logging subsystem. Can be empty.
//     severities: mask of logging severities
//     verbosity: logging verbosity
func AddConsoleLoggerStderr(
	name string,
	subsystem Subsystem,
	severities SeverityMask,
	verbosity Verbosity,
) {
	logger := NewConsoleLogger(os.Stderr)
	AddLogger(name, subsystem, severities, verbosity, logger)
}

// Add a Apache logger
//
// Parameters:
//     name: ID of the logger
//     subsystem: logging subsystem. Can be empty.
//     severities: mask of logging severities
//     verbosity: logging verbosity
//     file: path to the logging file
//     sync: flush all message immediately
func AddApacheLogger(
	name string,
	subsystem Subsystem,
	severities SeverityMask,
	verbosity Verbosity,
	file string,
	sync bool,
) {
	f := NewSimpleFile(file, sync)
	logger := NewApacheLogger(f)
	AddLogger(name, subsystem, severities, verbosity, logger)
}

// Add a Apache logger
//
// Parameters:
//     name: ID of the logger
//     subsystem: logging subsystem. Can be empty.
//     severities: mask of logging severities
//     verbosity: logging verbosity
//     pattern: pattern of names the log files
//     sync: flush all message immediately
func AddPatternApacheLogger(
	name string,
	subsystem Subsystem,
	severities SeverityMask,
	verbosity Verbosity,
	pattern string,
	sync bool,
) {
	f := NewPatternFile(timeSource, pattern, sync)
	logger := NewApacheLogger(f)
	AddLogger(name, subsystem, severities, verbosity, logger)
}

// Log a logging object into the global log
//
// Parameters:
//     subsystem: logging subsystem
//     severity: logging severity
//     verbosity: logging verbosity
//     object: logging object
func LogObject(
	subsystem Subsystem,
	severity Severity,
	verbosity Verbosity,
	object interface{},
) {
	DispatcherLogObject(globalLog, subsystem, severity, verbosity, object)
}

// Log a text message
//
// Parameters:
//     subsystem: logging subsystem
//     severity: logging severity
//     verbosity: logging verbosity
//     message: the message
func LogMessage(
	subsystem Subsystem,
	severity Severity,
	verbosity Verbosity,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, severity, verbosity, message)
}

// Log a formatted text message
//
// Parameters:
//     subsystem: logging subsystem
//     severity: logging severity
//     verbosity: logging verbosity
//     format: printf-like format of the message
//     args: arguments of the message
func LogMessagef(
	subsystem Subsystem,
	severity Severity,
	verbosity Verbosity,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(
		globalLog, subsystem, severity, verbosity, format, args...)
}

// There are a set of convenient logging functions. Their names follow
// the pattern:
//     <severity><verbosity>[f][s]
//
//     f .... the message is formatted
//     s .... a subsystem is specified

/* -- critical errors */
func Critical1(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Critical, 1, message)
}

func Critical1s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Critical, 1, message)
}

func Critical1f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Critical, 1, format, args...)
}

func Critical1fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Critical, 1, format, args...)
}

/* -- errors */
func Error1(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Error, 1, message)
}

func Error1s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Error, 1, message)
}

func Error1f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Error, 1, format, args...)
}

func Error1fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Error, 1, format, args...)
}

func Error2(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Error, 2, message)
}

func Error2s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Error, 2, message)
}

func Error2f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Error, 2, format, args...)
}

func Error2fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Error, 2, format, args...)
}

func Error3(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Error, 3, message)
}

func Error3s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Error, 3, message)
}

func Error3f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Error, 3, format, args...)
}

func Error3fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Error, 3, format, args...)
}

func Error4(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Error, 4, message)
}

func Error4s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Error, 4, message)
}

func Error4f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Error, 4, format, args...)
}

func Error4fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Error, 4, format, args...)
}

/* -- warnings */

func Warning1(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Warning, 1, message)
}

func Warning1s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Warning, 1, message)
}

func Warning1f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Warning, 1, format, args...)
}

func Warning1fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Warning, 1, format, args...)
}

func Warning2(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Warning, 2, message)
}

func Warning2s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Warning, 2, message)
}

func Warning2f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Warning, 2, format, args...)
}

func Warning2fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Warning, 2, format, args...)
}

func Warning3(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Warning, 3, message)
}

func Warning3s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Warning, 3, message)
}

func Warning3f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Warning, 3, format, args...)
}

func Warning3fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Warning, 3, format, args...)
}

func Warning4(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Warning, 4, message)
}

func Warning4s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Warning, 4, message)
}

func Warning4f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Warning, 4, format, args...)
}

func Warning4fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Warning, 4, format, args...)
}

/* -- info messages */
func Info1(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Info, 1, message)
}

func Info1s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Info, 1, message)
}

func Info1f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Info, 1, format, args...)
}

func Info1fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Info, 1, format, args...)
}

func Info2(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Info, 2, message)
}

func Info2s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Info, 2, message)
}

func Info2f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Info, 2, format, args...)
}

func Info2fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Info, 2, format, args...)
}

func Info3(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Info, 3, message)
}

func Info3s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Info, 3, message)
}

func Info3f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Info, 3, format, args...)
}

func Info3fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Info, 3, format, args...)
}

func Info4(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Info, 4, message)
}

func Info4s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Info, 4, message)
}

func Info4f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Info, 4, format, args...)
}

func Info4fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Info, 4, format, args...)
}

func Info5(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Info, 5, message)
}

func Info5s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Info, 5, message)
}

func Info5f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Info, 5, format, args...)
}

func Info5fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Info, 5, format, args...)
}

/* -- debug messages */

func Debug3(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Debug, 3, message)
}

func Debug3s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Debug, 3, message)
}

func Debug3f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Debug, 3, format, args...)
}

func Debug3fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Debug, 3, format, args...)
}

func Debug4(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Debug, 4, message)
}

func Debug4s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Debug, 4, message)
}

func Debug4f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Debug, 4, format, args...)
}

func Debug4fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Debug, 4, format, args...)
}

func Debug5(
	message string,
) {
	DispatcherLogMessage(globalLog, "", Debug, 5, message)
}

func Debug5s(
	subsystem Subsystem,
	message string,
) {
	DispatcherLogMessage(globalLog, subsystem, Debug, 5, message)
}

func Debug5f(
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, "", Debug, 5, format, args...)
}

func Debug5fs(
	subsystem Subsystem,
	format string,
	args ...interface{},
) {
	DispatcherLogMessagef(globalLog, subsystem, Debug, 5, format, args...)
}
