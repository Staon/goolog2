package goolog2

import (
	"runtime"
	"testing"
)

type testLogger struct {
	touched bool
}

func (this *testLogger) Destroy() {
	/* -- nothing to do */
}

func (this *testLogger) LogObject(
	system string,
	subsystem Subsystem,
	severity Severity,
	verbosity Verbosity,
	object interface{},
) {
	this.touched = true
}

func (this *testLogger) Check(
	expected bool,
) bool {
	retval := expected == this.touched
	this.touched = false
	return retval
}

const loggerCount = 12

func checkResult(
	t *testing.T,
	logger [loggerCount]*testLogger,
	expected string,
) {
	for i := 0; i < loggerCount; i++ {
		if !logger[i].Check(expected[i] != '0') {
			_, file, line, _ := runtime.Caller(1)
			t.Errorf("[%s:%d]: logger %d failed", file, line, i)
		}
	}
}

func TestLogDispatcher(t *testing.T) {
	Init("testlog")

	/* -- create loggers */
	var logger [loggerCount]*testLogger
	for i := 0; i < loggerCount; i++ {
		logger[i] = &testLogger{}
	}

	/* -- register the loggers */
	AddLogger("black hole", "", MaskAll, 0, logger[0])
	AddLogger("level 1", "", MaskAll, 1, logger[1])
	AddLogger("level 2", "", MaskAll, 2, logger[2])
	AddLogger("level 3", "", MaskAll, 3, logger[3])
	AddLogger("level 4", "", MaskAll, 4, logger[4])
	AddLogger("level 5", "", MaskAll, 5, logger[5])
	AddLogger("critical", "", MaskCritical, 5, logger[6])
	AddLogger("error", "", MaskError, 5, logger[7])
	AddLogger("warning", "", MaskWarning, 5, logger[8])
	AddLogger("info", "", MaskInfo, 5, logger[9])
	AddLogger("debug", "", MaskDebug, 5, logger[10])
	AddLogger("subsystem", "test", MaskAll, 5, logger[11])

	/* -- wrong subsystem */
	Critical1s("wrong", "")
	checkResult(t, logger, "011111100000")

	/* -- critical message */
	Critical1("")
	checkResult(t, logger, "011111100000")
	Critical1f("")
	checkResult(t, logger, "011111100000")
	Critical1s("test", "")
	checkResult(t, logger, "011111100001")
	Critical1fs("test", "")
	checkResult(t, logger, "011111100001")

	/* -- error messages */
	Error1("")
	checkResult(t, logger, "011111010000")
	Error1f("")
	checkResult(t, logger, "011111010000")
	Error1s("test", "")
	checkResult(t, logger, "011111010001")
	Error1fs("test", "")
	checkResult(t, logger, "011111010001")
	Error2("")
	checkResult(t, logger, "001111010000")
	Error2f("")
	checkResult(t, logger, "001111010000")
	Error2s("test", "")
	checkResult(t, logger, "001111010001")
	Error2fs("test", "")
	checkResult(t, logger, "001111010001")
	Error3("")
	checkResult(t, logger, "000111010000")
	Error3f("")
	checkResult(t, logger, "000111010000")
	Error3s("test", "")
	checkResult(t, logger, "000111010001")
	Error3fs("test", "")
	checkResult(t, logger, "000111010001")
	Error4("")
	checkResult(t, logger, "000011010000")
	Error4f("")
	checkResult(t, logger, "000011010000")
	Error4s("test", "")
	checkResult(t, logger, "000011010001")
	Error4fs("test", "")
	checkResult(t, logger, "000011010001")

	/* -- warning messages */
	Warning1("")
	checkResult(t, logger, "011111001000")
	Warning1f("")
	checkResult(t, logger, "011111001000")
	Warning1s("test", "")
	checkResult(t, logger, "011111001001")
	Warning1fs("test", "")
	checkResult(t, logger, "011111001001")
	Warning2("")
	checkResult(t, logger, "001111001000")
	Warning2f("")
	checkResult(t, logger, "001111001000")
	Warning2s("test", "")
	checkResult(t, logger, "001111001001")
	Warning2fs("test", "")
	checkResult(t, logger, "001111001001")
	Warning3("")
	checkResult(t, logger, "000111001000")
	Warning3f("")
	checkResult(t, logger, "000111001000")
	Warning3s("test", "")
	checkResult(t, logger, "000111001001")
	Warning3fs("test", "")
	checkResult(t, logger, "000111001001")
	Warning4("")
	checkResult(t, logger, "000011001000")
	Warning4f("")
	checkResult(t, logger, "000011001000")
	Warning4s("test", "")
	checkResult(t, logger, "000011001001")
	Warning4fs("test", "")
	checkResult(t, logger, "000011001001")

	/* -- warning messages */
	Info1("")
	checkResult(t, logger, "011111000100")
	Info1f("")
	checkResult(t, logger, "011111000100")
	Info1s("test", "")
	checkResult(t, logger, "011111000101")
	Info1fs("test", "")
	checkResult(t, logger, "011111000101")
	Info2("")
	checkResult(t, logger, "001111000100")
	Info2f("")
	checkResult(t, logger, "001111000100")
	Info2s("test", "")
	checkResult(t, logger, "001111000101")
	Info2fs("test", "")
	checkResult(t, logger, "001111000101")
	Info3("")
	checkResult(t, logger, "000111000100")
	Info3f("")
	checkResult(t, logger, "000111000100")
	Info3s("test", "")
	checkResult(t, logger, "000111000101")
	Info3fs("test", "")
	checkResult(t, logger, "000111000101")
	Info4("")
	checkResult(t, logger, "000011000100")
	Info4f("")
	checkResult(t, logger, "000011000100")
	Info4s("test", "")
	checkResult(t, logger, "000011000101")
	Info4fs("test", "")
	checkResult(t, logger, "000011000101")
	Info5("")
	checkResult(t, logger, "000001000100")
	Info5f("")
	checkResult(t, logger, "000001000100")
	Info5s("test", "")
	checkResult(t, logger, "000001000101")
	Info5fs("test", "")
	checkResult(t, logger, "000001000101")

	/* -- debug messages */
	Debug3("")
	checkResult(t, logger, "000111000010")
	Debug3f("")
	checkResult(t, logger, "000111000010")
	Debug3s("test", "")
	checkResult(t, logger, "000111000011")
	Debug3fs("test", "")
	checkResult(t, logger, "000111000011")
	Debug4("")
	checkResult(t, logger, "000011000010")
	Debug4f("")
	checkResult(t, logger, "000011000010")
	Debug4s("test", "")
	checkResult(t, logger, "000011000011")
	Debug4fs("test", "")
	checkResult(t, logger, "000011000011")
	Debug5("")
	checkResult(t, logger, "000001000010")
	Debug5f("")
	checkResult(t, logger, "000001000010")
	Debug5s("test", "")
	checkResult(t, logger, "000001000011")
	Debug5fs("test", "")
	checkResult(t, logger, "000001000011")

	Destroy()
}
