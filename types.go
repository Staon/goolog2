package goolog2

import ()

// Severity of a log message
type Severity uint32

const (
	Critical Severity = 1 << iota
	Error
	Warning
	Info
	Debug
)

func (this Severity) Code() string {
	switch this {
	case Critical:
		return "CRITICAL"
	case Error:
		return "ERROR"
	case Warning:
		return "WARNING"
	case Info:
		return "INFO"
	case Debug:
		return "DEBUG"
	default:
		panic("invalid severity code")
	}
}

// Mask of severities
type SeverityMask uint32

const (
	MaskCritical SeverityMask = SeverityMask(Critical)
	MaskError    SeverityMask = SeverityMask(Error)
	MaskWarning  SeverityMask = SeverityMask(Warning)
	MaskInfo     SeverityMask = SeverityMask(Info)
	MaskDebug    SeverityMask = SeverityMask(Debug)
	MaskStd      SeverityMask = SeverityMask(Critical | Error | Warning | Info)
	MaskAll      SeverityMask = SeverityMask(Critical | Error | Warning | Info | Debug)
)

// Verbosity of a log message
//
// 0 means no logging, a higher number means a higher verbosity level
type Verbosity uint32
