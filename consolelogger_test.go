package goolog2_test

import (
	. "goolog2"
	"testing"
)

func TestConsoleLogger(t *testing.T) {
	Init("testlog")
	AddConsoleLoggerStderr("console", "", MaskAll, 5)

	/* -- I cannot check the output into the console - the user must
	   do that by yes. */
	Critical1("critical message")
	Error2("error message")
	Warning3("warning message")
	Info4("info message")
	Debug5("debug message")

	Destroy()
}
