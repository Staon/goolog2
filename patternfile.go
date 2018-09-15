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
	timesrc     TimeSource
	pattern     string
	sync        bool
	nextCheck   int64
	currName    string
	currFile    *os.File
	currWriter  FileWriter
	lineMutex   sync.Mutex
	rotateMutex sync.Mutex
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
func NewPatternFile(
	timesrc TimeSource,
	pattern string,
	sync bool,
) FileHolder {
	return &patternFile{
		timesrc: timesrc,
		pattern: pattern,
		sync:    sync,
	}
}

func (this *patternFile) AccessWriter(
	functor func(writer FileWriter),
) {
	/* -- check whether new filename should be tried */
	now := this.timesrc.Now()
	checkTime := atomic.LoadInt64(&this.nextCheck)
	if now.Unix() >= checkTime {
		/* -- try to generate new filename */
		this.checkNewFile(now)
	}

	/* -- log the line */
	this.lineMutex.Lock()
	defer this.lineMutex.Unlock()
	if this.currFile != nil {
		functor(this.currWriter)
		if this.sync {
			this.currFile.Sync()
		}
	}
}

func (this *patternFile) Close() {

}

func (this *patternFile) checkNewFile(
	now time.Time,
) {
	/* -- synchronize threads which try to rotate the log file */
	this.rotateMutex.Lock()
	defer this.rotateMutex.Unlock()

	/* -- another thread could already rotate the file */
	if now.Unix() >= this.nextCheck {
		/* -- compute and set time of next check */
		nextCheck := now.Unix() + checkInterval
		atomic.StoreInt64(&this.nextCheck, nextCheck)

		/* -- generate new filename */
		newName := this.generateFilename(now)

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
			oldFile := this.currFile
			this.currFile = newFile
			this.currWriter = newSimpleFileWriter(newFile)
			this.lineMutex.Unlock()

			/* -- close the old file */
			if oldFile != nil {
				oldFile.Close()
			}
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
