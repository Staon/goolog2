package goolog2_test

import (
	"os"
	"testing"
	"time"

	. "github.com/Staon/goolog2"
)

func TestRotatableFile(t *testing.T) {
	const long_message_base = " very very very very very very long long long long long long long long long long long long long long long long long long"
	const long_message = long_message_base + long_message_base + " message"
	defer func() {
		rotatableLogClean("")
		Destroy()
	}()

	timesrc := rotatableLogInit(t, "")

	/* -- Create the logger */
	AddRotatableFileLogger(
		"file", "", MaskAll, 5, "rotatable.log", false, 200, 2*time.Minute)

	Error1("First error")
	testRotatableLogExists(t, "rotatable", false, false, false, "0 minutes")

	/* -- Shift the time by half of checking interval. */
	Error1("Second error")
	timesrc.ShiftTime(time.Minute)
	testRotatableLogExists(t, "rotatable", false, false, false, "1 minute")

	/* -- Shift the time by more than checking interval. */
	//    The log file is too short to rotate.
	Error1("Third error")
	timesrc.ShiftTime(3 * time.Minute)
	testRotatableLogExists(t, "rotatable", false, false, false, "4 minutes")

	/* -- Shift the time by more than checking interval */
	//    The log file is long enough to rotate.
	Error1("First " + long_message)
	timesrc.ShiftTime(3 * time.Minute)
	testRotatableLogExists(t, "rotatable", true, false, false, "7 minutes")

	/* -- Shift the time by half of checking interval. */
	//    The log file is long enough to rotate but the time to check file it not yet.
	Error1("Second " + long_message)
	timesrc.ShiftTime(time.Minute)
	testRotatableLogExists(t, "rotatable", true, false, false, "8 minutes")

	/* -- Shift the time by checking interval. */
	//    The log file is long enough to rotate.
	Error1("Short message")
	timesrc.ShiftTime(3 * time.Minute)
	testRotatableLogExists(t, "rotatable", true, true, false, "11 minutes")
}

func TestRotatableRecovery(t *testing.T) {
	const long_message_base = " very very very very very very long long long long long long long long long long long long long long long long long long"
	const long_message = long_message_base + long_message_base + " message"

	defer func() {
		rotatableLogClean("")
		Destroy()
	}()

	timesrc := rotatableLogInit(t, "")
	f, _ := os.Create("rotatable.log.2")
	f.Close()

	AddRotatableFileLogger(
		"file", "", MaskAll, 5, "rotatable.log", false, 200, 2*time.Minute)
	/* --  log exists, log.1 doesn't exist, log.2 exists (a hole in the sequence) --*/
	//     A simulation of the state "a renaming log to log.1 has failed".
	testRotatableLogExists(t, "rotatable", false, true, false, "0 minutes")
	/* -- Create log.1 --*/
	Error1("First " + long_message)
	Error1("First short error")
	timesrc.ShiftTime(3 * time.Minute)
	testRotatableLogExists(t, "rotatable", true, true, false, "3 minutes")
	/* -- Now is the sequence without holes. It works as usual. --*/
	Error1("Second " + long_message)
	Error1("Second short error")
	timesrc.ShiftTime(3 * time.Minute)
	testRotatableLogExists(t, "rotatable", true, true, true, "6 minutes")
}

func TestMoreRotatableLogs(t *testing.T) {
	const long_message_base = " very very very very very very long long long long long long long long long long long long long long long long long long"
	const long_message = long_message_base + long_message_base + " message"

	suffixes := []string{"-2m", "-5m", "-8m"}
	defer func() {
		rotatableLogClean(suffixes...)
		os.Remove("rotatable-2m.log.4")
		os.Remove("rotatable-2m.log.5")
		Destroy()
	}()

	/* -- Create loggers, every has different check interval  -- */
	timesrc := rotatableLogInit(t, suffixes...)
	for i, s := range suffixes {
		AddRotatableFileLogger(
			"file"+s, "", MaskAll, 5, "rotatable"+s+".log", false, 200, time.Duration(2+i*3)*time.Minute)
	}
	/* -- Initial short error (none rotation). -- */
	Error1("First error")
	for _, s := range suffixes {
		testRotatableLogExists(t, "rotatable"+s, false, false, false, "0 minutes")
	}
	/* -- Shift the time by 3 minutes - the time to the rotate: */
	//    first - now, second - 2:00, third - 5:00
	Error1("First " + long_message)
	timesrc.ShiftTime(3 * time.Minute)
	for i, s := range suffixes {
		testRotatableLogExists(t, "rotatable"+s, i < 1, false, false, "3 minutes")
	}
	/* -- Shift the time the time by 3 minutes - the time to the rotate: */
	//    first - now, second - now, third - 2:00
	Error1("Second " + long_message)
	timesrc.ShiftTime(3 * time.Minute)
	for i, s := range suffixes {
		testRotatableLogExists(t, "rotatable"+s, i < 2, i < 1, false, "6 minutes")
	}
	/* -- Shift the time by 4 minutes - the time to the rotate: */
	//    first - now, second - 1:00, third - now
	Error1("Third " + long_message)
	timesrc.ShiftTime(4 * time.Minute)
	for i, s := range suffixes {
		testRotatableLogExists(t, "rotatable"+s, true, i < 1, i < 1, "10 minutes")
	}
	/* -- Shift the time by one and half minutes - the time to the rotate: */
	//   first - 0:30, second - now, third - 7:30
	Error1("Fourth " + long_message)
	timesrc.ShiftTime(90 * time.Second)
	for i, s := range suffixes {
		testRotatableLogExists(t, "rotatable"+s, true, i < 2, i < 1, "11 minutes 30 seconds")
	}

}

func rotatableLogInit(t *testing.T, suffixes ...string) *mockTimeSource {
	now, _ := time.Parse("2006-01-02T15:04:05", "2018-08-25T14:02:00")
	timesrc := &mockTimeSource{
		now: now,
	}
	rotatableLogClean(suffixes...)
	for _, s := range suffixes {
		// The function  testRotatableLogExists assumes that the main log file always exists.
		f, _ := os.Create("rotatable" + s + ".log")
		f.Close()
		testRotatableLogExists(t, "rotatable"+s, false, false, false, "initial state")
	}
	InitWithTimeSource("testlog", timesrc)
	return timesrc
}

func rotatableLogClean(suffixes ...string) {
	for _, s := range suffixes {
		os.Remove("rotatable" + s + ".log")
		os.Remove("rotatable" + s + ".log.1")
		os.Remove("rotatable" + s + ".log.2")
		os.Remove("rotatable" + s + ".log.3")
	}
}

func testRotatableLogExists(t *testing.T, fileBase string, shouldExist1 bool, shouldExist2 bool, shouldExist3 bool, messageId string) {
	testRotatableLogExistsImpl(t, fileBase+".log", true, messageId)
	testRotatableLogExistsImpl(t, fileBase+".log.1", shouldExist1, messageId)
	testRotatableLogExistsImpl(t, fileBase+".log.2", shouldExist2, messageId)
	testRotatableLogExistsImpl(t, fileBase+".log.3", shouldExist3, messageId)
}

func testRotatableLogExistsImpl(t *testing.T, fileName string, shouldExist bool, messageId string) {
	_, err := os.Stat(fileName)
	if err != nil && shouldExist {
		t.Fatalf("Error - %s: The file '%s' doesn't exists.", messageId, fileName)
	}
	if err == nil && !shouldExist {
		t.Fatalf("Error - %s: The file '%s' exists.", messageId, fileName)
	}
}
