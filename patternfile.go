package goolog2

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const checkInterval int64 = 30

type patternFile struct {
	pattern    string
	sync       bool
	currName   string
	currWriter FileWriter
	lineMutex  sync.Mutex
	refcount   int32
}

// Create new pattern file holder
//
// This implementation allows rotating of log files according
// specified pattern.
//
// The pattern is a string which can contain special sequences:
//     %Y ..... year (4 digits)
//     %m ..... month (01 - 12)
//     %d ..... day (01 - 31)
//     %H ..... hour (00 - 23)
//     %M ..... minute (00 - 59)
//     %% ..... %
//
// Parameters:
//     timesrc: time source
//     pattern: the filename pattern
//     sync: flush the file after every line
// Returns:
//     the new file holder
// Note: the reference counter is set to 1. You have to invoke Unref()
//     to clean up the holder.
func NewPatternFile(
	timesrc TimeSource,
	pattern string,
	sync bool,
) RotatableFileHolder {
	holder := patternFile{
		pattern:  pattern,
		sync:     sync,
		refcount: 1,
	}
	holder.Rotate(timesrc)
	return &holder
}

func (this *patternFile) AccessWriter(
	functor func(writer FileWriter),
) {
	/* -- log the line */
	this.lineMutex.Lock()
	defer this.lineMutex.Unlock()
	if this.currWriter != nil {
		functor(this.currWriter)
		if this.sync {
			this.currWriter.Sync()
		}
	}
}

func (this *patternFile) Ref() FileHolder {
	atomic.AddInt32(&this.refcount, 1)
	return this
}

func (this *patternFile) Unref() {
	refcount := atomic.AddInt32(&this.refcount, -1)
	if refcount == 0 && this.currWriter != nil {
		this.currWriter.Close()
		this.currWriter = nil
	}
}

// See LogRotator interface
func (this *patternFile) NeedRotate(
	timesrc TimeSource,
) bool {
	return true
}

// See LogRotator interface
func (this *patternFile) Rotate(
	timesrc TimeSource,
) {
	/* -- generate new filename */
	newName := this.generateFilename(timesrc.Now())

	/* -- if the filename differs switch the files */
	if newName != this.currName {

		/* -- No one other reads this value. So I can store it
		without any synchronization */
		this.currName = newName

		/* -- open the new file */
		newFile, err := os.OpenFile(
			newName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			newFile = nil
		}

		/* -- switch the file */
		this.lineMutex.Lock()
		oldWriter := this.currWriter
		this.currWriter = newSimpleFileWriter(newFile, true)
		this.lineMutex.Unlock()

		/* -- close the old file */
		if oldWriter != nil {
			oldWriter.Close()
		}
	}
}

func (this *patternFile) generateFilename(
	now time.Time,
) string {
	type stateCode int
	const (
		INIT stateCode = iota
		FMT
	)

	pattern := this.pattern
	builder := &strings.Builder{}
	state := INIT
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		switch state {
		case INIT:
			if c == '%' {
				state = FMT
			} else {
				builder.WriteByte(c)
			}
		case FMT:
			switch c {
			case 'Y':
				fmt.Fprintf(builder, "%04d", now.Year())
			case 'm':
				fmt.Fprintf(builder, "%02d", now.Month())
			case 'd':
				fmt.Fprintf(builder, "%02d", now.Day())
			case 'H':
				fmt.Fprintf(builder, "%02d", now.Hour())
			case 'M':
				fmt.Fprintf(builder, "%02d", now.Minute())
			default:
				builder.WriteByte(c)
			}
			state = INIT
		}
	}

	return builder.String()
}

// See LogRotator interface
func (this *patternFile) GetNextCheckTime(
	timesrc TimeSource,
) time.Time {
	/* -- compute and set time of next check */
	now := timesrc.Now()
	nextCheck := now.Add(time.Duration(checkInterval) * time.Second)
	return nextCheck
}
