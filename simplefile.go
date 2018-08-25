package goolog2

import (
	"io"
	"os"
	"sync"
)

type simpleFile struct {
	writer *os.File
	mutex  sync.Mutex
	sync   bool
}

func NewSimpleFile(
	filepath string,
	sync bool,
) FileHolder {
	holder := &simpleFile{
		sync: sync,
	}

	// I ignore the error here - if the file cannot be opened, the logging
	// just simply doesn't work. However, if I cannot open the logging
	// file, I cannot report the error.
	holder.writer, _ = os.OpenFile(
		filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)

	return holder
}

func (this *simpleFile) AccessWriter(
	functor func(writer io.Writer),
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

func (this *simpleFile) Close() {
	if this.writer != nil {
		this.writer.Close()
		this.writer = nil
	}
}
