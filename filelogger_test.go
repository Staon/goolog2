package goolog2_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
	"time"

	. "github.com/Staon/goolog2"
)

type mockTimeSource struct {
	now time.Time
}

func (this *mockTimeSource) Now() time.Time {
	return this.now
}

// Change the time and wait to rotators actions scheduled for this time or earlier.
// Paramaters:
//      newTime - string in the format "2006-01-02T15:04:05"
// Returns:
//      New mocked time
func (this *mockTimeSource) SetTime(newTime string) time.Time {
	now, _ := time.Parse("2006-01-02T15:04:05", newTime)
	if this.now != now {
		this.now = now
		AfterChangeMockedTime(true)
	}
	return this.now
}

// Shift the time and wait to rotators actions scheduled for this time or earlier.
// Paramaters:
//      duration
// Returns:
//      New mocked time
func (this *mockTimeSource) ShiftTime(duration time.Duration) time.Time {
	if duration > 0 {
		this.now = this.now.Add(duration)
		AfterChangeMockedTime(true)
	}
	return this.now
}

func TestFileLogger(t *testing.T) {
	logfile := "log.log"
	os.Remove(logfile)

	now, _ := time.Parse("2006-01-02T15:04:05", "2018-08-25T14:02:27")
	timesrc := &mockTimeSource{
		now: now,
	}

	InitWithTimeSource("testlog", timesrc)
	AddFileLogger("file", "", MaskAll, 5, logfile, false)

	/* -- critical messages */
	Critical1("critical error")
	Critical1f("critical %s", "error")
	Critical1s("test", "s critical error")
	Critical1fs("test", "s critical %s", "error")

	/* -- error messages */
	Error1("error")
	Error1f("%s", "error")
	Error1s("test", "s error")
	Error1fs("test", "s %s", "error")
	Error2("error")
	Error2f("%s", "error")
	Error2s("test", "s error")
	Error2fs("test", "s %s", "error")
	Error3("error")
	Error3f("%s", "error")
	Error3s("test", "s error")
	Error3fs("test", "s %s", "error")
	Error4("error")
	Error4f("%s", "error")
	Error4s("test", "s error")
	Error4fs("test", "s %s", "error")

	/* -- warning messages */
	Warning1("warning")
	Warning1f("%s", "warning")
	Warning1s("test", "s warning")
	Warning1fs("test", "s %s", "warning")
	Warning2("warning")
	Warning2f("%s", "warning")
	Warning2s("test", "s warning")
	Warning2fs("test", "s %s", "warning")
	Warning3("warning")
	Warning3f("%s", "warning")
	Warning3s("test", "s warning")
	Warning3fs("test", "s %s", "warning")
	Warning4("warning")
	Warning4f("%s", "warning")
	Warning4s("test", "s warning")
	Warning4fs("test", "s %s", "warning")

	/* -- info messages */
	Info1("info")
	Info1f("%s", "info")
	Info1s("test", "s info")
	Info1fs("test", "s %s", "info")
	Info2("info")
	Info2f("%s", "info")
	Info2s("test", "s info")
	Info2fs("test", "s %s", "info")
	Info3("info")
	Info3f("%s", "info")
	Info3s("test", "s info")
	Info3fs("test", "s %s", "info")
	Info4("info")
	Info4f("%s", "info")
	Info4s("test", "s info")
	Info4fs("test", "s %s", "info")
	Info5("info")
	Info5f("%s", "info")
	Info5s("test", "s info")
	Info5fs("test", "s %s", "info")

	/* -- debug messages */
	Debug3("debug")
	Debug3f("%s", "debug")
	Debug3s("test", "s debug")
	Debug3fs("test", "s %s", "debug")
	Debug4("debug")
	Debug4f("%s", "debug")
	Debug4s("test", "s debug")
	Debug4fs("test", "s %s", "debug")
	Debug5("debug")
	Debug5f("%s", "debug")
	Debug5s("test", "s debug")
	Debug5fs("test", "s %s", "debug")

	Destroy()

	/* -- check the log */
	expected, err1 := ioutil.ReadFile("filelogger.log")
	current, err2 := ioutil.ReadFile(logfile)
	if err1 != nil || err2 != nil || !bytes.Equal(expected, current) {
		t.Errorf("generated log file is different!")
	}
}
