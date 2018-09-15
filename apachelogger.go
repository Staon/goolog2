package goolog2

import (
	"fmt"
)

type apacheLogger struct {
	file FileHolder
}

// Create new apache logger
//
// The Apache logger expects the Apache logging objects. The output
// format follows the combined Apache logging format.
//
// Parameters:
//     file: a file holder used for the output
// Returns:
//     the logger
func NewApacheLogger(
	file FileHolder,
) Logger {
	return &apacheLogger{
		file: file,
	}
}

func (this *apacheLogger) Destroy() {
	this.file.Close()
}

func (this *apacheLogger) LogObject(
	system string,
	subsystem Subsystem,
	severity Severity,
	verbosity Verbosity,
	object interface{},
) {
	/* -- the logger supports only line objects */
	apache, ok := object.(ApacheObject)
	if !ok {
		return
	}

	/* -- get the values */
	remoteHost, identity, user, requestTime, method, resource, protocol,
		status, length, referer, agent := apache.GetValues()
	if remoteHost == "" {
		remoteHost = "-"
	}
	if identity == "" {
		identity = "-"
	}
	if user == "" {
		user = "-"
	}
	if referer == "" {
		referer = "-"
	}
	if agent == "" {
		agent = "-"
	}

	/* -- write the message */
	this.file.AccessWriter(func(writer FileWriter) {
		fmt.Fprintf(
			writer,
			"%s %s %s [%s] \"%s %s %s\" %03d %d \"%s\" \"%s\"\n",
			remoteHost,
			identity,
			user,
			requestTime.Format("02/Jan/2006:15:04:05 -0700"),
			method,
			resource,
			protocol,
			status,
			length,
			referer,
			agent)
	})
}
