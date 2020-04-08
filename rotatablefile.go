package goolog2

import (
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type rotatableFile struct {
	filePath      string
	maxSize       int64
	writer        FileWriter
	checkInterval time.Duration
	mutex         sync.Mutex
	sync          bool
	refcount      int32
}

// Create new rotatable file holder.
//
// Parameters:
//     filepath: name of the logging file
//     sync: if true, the stream is flushed after every message
//     maxSize:  Make log rotation if log size is bigger than maxSize.
//     checkInterval: time interval to check the log size; usually minutes or tens of minutes
// Returns:
//     the new rotatable file holder
// Note: the reference counter is set to 1. You have to invoke Unref()
//     to clean up the holder.
func NewRotatableFile(filePath string, sync bool, maxSize int64, checkInterval time.Duration) RotatableFileHolder {
	holder := rotatableFile{
		filePath:      filePath,
		maxSize:       maxSize,
		checkInterval: checkInterval,
		sync:          sync,
		refcount:      1,
	}

	// I ignore the error here - if the file cannot be opened, the logging
	// just simply doesn't work. However, if I cannot open the logging
	// file, I cannot report the error.
	file, _ := os.OpenFile(
		filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	holder.writer = newSimpleFileWriter(file, true)
	return &holder
}

func (this *rotatableFile) AccessWriter(
	functor func(writer FileWriter),
) {
	if this.writer != nil {
		/* --  The lock avoids inter-mixing of logging lines */
		this.mutex.Lock()
		defer this.mutex.Unlock()
		functor(this.writer)
		if this.sync {
			this.writer.Sync()
		}
	}
}

func (this *rotatableFile) Ref() FileHolder {
	atomic.AddInt32(&this.refcount, 1)
	return this
}

func (this *rotatableFile) Unref() {
	refcount := atomic.AddInt32(&this.refcount, -1)
	if refcount == 0 {
		if this.writer != nil {
			this.writer.Close()
			this.writer = nil
		}
	}
}

// See LogRotator interface
func (this *rotatableFile) NeedRotate(timesrc TimeSource) bool {
	if this.maxSize == 0 {
		return false
	}
	this.writer.Sync()
	fileInfo := this.writer.Stat()
	return fileInfo != nil && fileInfo.Size() > this.maxSize
}

// See LogRotator interface
func (this *rotatableFile) Rotate(timesrc TimeSource) {
	// A error in this part is not fatal. It will be recovered in next successfull Rotate().
	i := 0
	var err error
	for err == nil {
		i++
		_, err = os.Stat(this.filePath + "." + strconv.Itoa(i))
	}
	for ; i > 1; i-- {
		source := this.filePath + "." + strconv.Itoa(i-1)
		target := this.filePath + "." + strconv.Itoa(i)
		if err := os.Rename(source, target); err != nil {
			return
		}
	}
	// rename current file
	this.mutex.Lock() // stop writing to file
	defer this.mutex.Unlock()
	this.writer.Close()
	os.Rename(this.filePath, this.filePath+".1")
	file, _ := os.OpenFile(
		this.filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	this.writer = newSimpleFileWriter(file, true)
}

// See LogRotator interface
func (this *rotatableFile) GetNextCheckTime(timesrc TimeSource) time.Time {
	return timesrc.Now().Add(this.checkInterval)
}
