package goolog2

import (
	"os"

	atty "github.com/mattn/go-isatty"
	"github.com/xo/terminfo"
)

type simpleFileWriter struct {
	file  *os.File
	owner bool
	tinfo *terminfo.Terminfo
}

func newSimpleFileWriter(
	file *os.File,
	owner bool,
) FileWriter {
	writer := &simpleFileWriter{
		file:  file,
		owner: owner,
	}

	/* -- get the terminfo object */
	if isatty := atty.IsTerminal(file.Fd()); isatty {
		tinfo, err := terminfo.LoadFromEnv()
		if err == nil {
			writer.tinfo = tinfo
		}
	}

	return writer
}

func (this *simpleFileWriter) Close() error {
	if !this.owner || this.file == nil {
		return nil
	}
	err := this.file.Close()
	this.file = nil
	return err
}

func (this *simpleFileWriter) Stat() os.FileInfo {
	if this.file == nil {
		return nil
	}
	stat, err := this.file.Stat()
	if err != nil {
		return nil
	}
	return stat
}

func (this *simpleFileWriter) Sync() {
	if this.file != nil {
		this.file.Sync()
	}
}

func (this *simpleFileWriter) Write(
	p []byte,
) (int, error) {
	if this.file == nil {
		return 0, os.ErrClosed
	}
	return this.file.Write(p)
}

func (this *simpleFileWriter) ChangeColor(
	color Color,
) {
	if this.tinfo != nil && this.file != nil {
		/* -- select color */
		var tcolor int
		switch color {
		case RED:
			tcolor = 1 /* -- red */
		case YELLOW:
			tcolor = 3 /* -- yellow */
		case BLUE:
			tcolor = 2 /* -- green */
		default:
			tcolor = -1
		}

		/* -- colorize the output */
		if tcolor > 0 {
			this.file.WriteString(
				this.tinfo.Printf(terminfo.SetAForeground, tcolor))
		}
	}
}

func (this *simpleFileWriter) ResetColor() {
	if this.tinfo != nil && this.file != nil {
		this.file.WriteString(
			this.tinfo.Printf(terminfo.ExitAttributeMode))
	}
}
