package goolog2

import (
	"fmt"
	"sync"
)

// Dispatch a log object into the logger objects
type LogDispatcher interface {
	// Destroy the distpatcher
	//
	// Warning: this method is not thread safe. Use it if you're certain
	//     all threads have stopped already.
	Destroy()

	// Log an object
	//
	// Parameters:
	//     subsystem: ID of logging subsystem
	//     severity: severity of the log message
	//     verbosity: verbosity of the log message
	//     object: an object to be logged
	LogObject(
		subsystem Subsystem,
		severity Severity,
		verbosity Verbosity,
		object interface{})

	// Add new logger
	//
	// Parameters:
	//     name: a unique name of the logger
	//     subsystem: ID of logging subsystem
	//     severity: maximal severity of the logger
	//     verbosity: maximal verbosity of the logger
	//     logger: the logger object
	AddLogger(
		name string,
		subsystem Subsystem,
		severities SeverityMask,
		verbosity Verbosity,
		logger Logger)
}

type logDispatcherRecord struct {
	subsystem  Subsystem
	severities SeverityMask
	verbosity  Verbosity
	logger     Logger
}

type logDispatcher struct {
	system  string
	loggers map[string]*logDispatcherRecord
	mutex   sync.RWMutex
}

// Create new log dispatcher
//
// Parameters:
//     system: an identifier shown in the logs
func NewLogDispatcher(
	system string,
) LogDispatcher {
	return &logDispatcher{
		system:  system,
		loggers: make(map[string]*logDispatcherRecord),
	}
}

func (this *logDispatcher) Destroy() {
	for _, record := range this.loggers {
		record.logger.Destroy()
	}
}

func (this *logDispatcher) LogObject(
	subsystem Subsystem,
	severity Severity,
	verbosity Verbosity,
	object interface{},
) {
	/* -- iterate loggers, find matching ones and log the object */
	this.mutex.RLock()
	defer this.mutex.RUnlock()
	for _, record := range this.loggers {
		if (record.subsystem == "" || record.subsystem == subsystem) &&
			(uint32(record.severities)&uint32(severity)) != 0 &&
			verbosity <= record.verbosity {
			/* -- the conditions match, log the object */
			record.logger.LogObject(
				this.system, subsystem, severity, verbosity, object)
		}
	}
}

func (this *logDispatcher) AddLogger(
	name string,
	subsystem Subsystem,
	severities SeverityMask,
	verbosity Verbosity,
	logger Logger,
) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.loggers[name] = &logDispatcherRecord{
		subsystem:  subsystem,
		severities: severities,
		verbosity:  verbosity,
		logger:     logger,
	}
}

// Log a logging object
func DispatcherLogObject(
	log LogDispatcher,
	subsystem Subsystem,
	severity Severity,
	verbosity Verbosity,
	object interface{},
) {
	log.LogObject(subsystem, severity, verbosity, object)
}

type simpleLogMessageObject struct {
	message string
}

func (this *simpleLogMessageObject) GetLogLine() string {
	return this.message
}

// Log a text message
func DispatcherLogMessage(
	log LogDispatcher,
	subsystem Subsystem,
	severity Severity,
	verbosity Verbosity,
	message string,
) {
	log.LogObject(
		subsystem,
		severity,
		verbosity,
		&simpleLogMessageObject{message})
}

type formattedLogMessageObject struct {
	format string
	args   []interface{}
}

func (this *formattedLogMessageObject) GetLogLine() string {
	return fmt.Sprintf(this.format, this.args...)
}

// Log a formatted text message
func DispatcherLogMessagef(
	log LogDispatcher,
	subsystem Subsystem,
	severity Severity,
	verbosity Verbosity,
	format string,
	args ...interface{},
) {
	log.LogObject(
		subsystem,
		severity,
		verbosity,
		&formattedLogMessageObject{format, args})
}
