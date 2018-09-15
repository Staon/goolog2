package goolog2

import ()

type fileLogger struct {
	timesrc   TimeSource
	file      FileHolder
	formatter LineFormatter
}

// Create new file logger
//
// Parameters:
//     timesrc: a timesource
//     file: a file holder used for the output
//     formatter: a line formatter
// Returns:
//     the logger
func NewFileLogger(
	timesrc TimeSource,
	file FileHolder,
	formatter LineFormatter,
) Logger {
	return &fileLogger{
		timesrc:   timesrc,
		file:      file,
		formatter: formatter,
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
	this.file.AccessWriter(func(writer FileWriter) {
		this.formatter.FormatLine(
			writer,
			this.timesrc.Now(),
			system,
			subsystem,
			severity,
			verbosity,
			line.GetLogLine())
	})
}
