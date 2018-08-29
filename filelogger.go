package goolog2

import (
	"fmt"
	"io"
)

type fileLogger struct {
	timesrc TimeSource
	file    FileHolder
}

// Create new file logger
//
// Parameters:
//     timesrc: a timesource
//     file: a file holder used for the output
// Returns:
//     the logger
func NewFileLogger(
	timesrc TimeSource,
	file FileHolder,
) Logger {
	return &fileLogger{
		timesrc: timesrc,
		file:    file,
	}
}

func (this *fileLogger) Destroy() {
	this.file.Close()
}

func (this *fileLogger) LogObject(
	system string,
	subsystem Subsystem,
	severity Severity,
	verbosity Verbosity,
	object interface{},
) {
	/* -- the logger supports only line objects */
	line, ok := object.(LineObject)
	if !ok {
		return
	}

	/* -- write the message */
	this.file.AccessWriter(func(writer io.Writer) {
		fmt.Fprintf(
			writer,
			"%s %s [%8s, %d] (%s): %s\n",
			system,
			this.timesrc.Now().Format("2006-01-02T15:04:05"),
			severity.Code(),
			verbosity,
			subsystem,
			line.GetLogLine())
	})
}
