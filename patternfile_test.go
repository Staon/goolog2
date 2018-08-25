package goolog2

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestPatternFile(t *testing.T) {
	now, _ := time.Parse("2006-01-02T15:04:05", "2018-08-25T14:02:00")
	timesrc := &mockTimeSource{
		now: now,
	}

	os.Remove("pattern2018-08-25-14:02.log")
	os.Remove("pattern2018-08-25-14:03.log")

	InitWithTimeSource("testlog", timesrc)

	/* -- create the logger */
	f := NewPatternFile(timesrc, "pattern%Y-%m-%d-%H:%M.log", false)
	logger := NewFileLogger(timesrc, f)
	AddLogger("file", "", MaskAll, 5, logger)

	Info1("First message")
	Error1("First error")

	/* -- move time after the checking interval */
	now, _ = time.Parse("2006-01-02T15:04:05", "2018-08-25T14:02:30")
	timesrc.now = now
	Info1("Second message")
	Error1("Second error")

	/* -- move time to force rotation of the logs */
	now, _ = time.Parse("2006-01-02T15:04:05", "2018-08-25T14:03:00")
	timesrc.now = now
	Info1("Third message")
	Error1("Third error")

	Destroy()

	/* -- check the logs */
	expected, err1 := ioutil.ReadFile("pattern1.log")
	current, err2 := ioutil.ReadFile("pattern2018-08-25-14:02.log")
	if err1 != nil || err2 != nil || !bytes.Equal(expected, current) {
		t.Errorf("first generated log file is different!")
	}
	expected, err1 = ioutil.ReadFile("pattern2.log")
	current, err2 = ioutil.ReadFile("pattern2018-08-25-14:03.log")
	if err1 != nil || err2 != nil || !bytes.Equal(expected, current) {
		t.Errorf("second generated log file is different!")
	}
}
