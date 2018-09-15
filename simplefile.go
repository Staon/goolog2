package goolog2

import (
	"os"
	"sync"
)

type simpleFile struct {
	owner  bool
	file   *os.File
	writer FileWriter
	mutex  sync.Mutex
	sync   bool
}

// Create new simple file holder
//
// Parameters:
//     filepath: name of the logging file
//     sync: if true, the stream is flushed after every message
func NewSimpleFile(
	filepath string,
	sync bool,
) FileHolder {
	holder := &simpleFile{
		owner: true,
		sync:  sync,
	}

	// I ignore the error here - if the file cannot be opened, the logging
	// just simply doesn't work. However, if I cannot open the logging
	// file, I cannot report the error.
	holder.file, _ = os.OpenFile(
		filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	holder.writer = newSimpleFileWriter(holder.file)

	return holder
}

// Create new simple file holder working over an opened file handle
//
// Parameters:
//     file: the opened file
//     sync: if true, the stream is flushed after every message
func NewSimpleFileHandle(
	file *os.File,
	sync bool,
) FileHolder {
	return &simpleFile{
		owner:  false,
		file:   file,
		writer: newSimpleFileWriter(file),
		sync:   sync,
	}
}

func (this *simpleFile) AccessWriter(
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

func (this *simpleFile) Close() {
	if this.writer != nil && this.owner {
		this.file.Close()
		this.file = nil
		this.writer = nil
	}
}
