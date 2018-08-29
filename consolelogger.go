package goolog2

import (
	"fmt"
	"os"
	"sync"

	atty "github.com/mattn/go-isatty"
	"github.com/xo/terminfo"
)

type consoleLogger struct {
	output *os.File
	tinfo  *terminfo.Terminfo
	isatty bool
	mutex  sync.Mutex
}

// Create new console logger
//
// Parameters:
//     output: an output stream
func NewConsoleLogger(
	output *os.File,
) Logger {
	logger := &consoleLogger{
		output: output,
	}

	/* -- get the terminfo object */
	if isatty := atty.IsTerminal(output.Fd()); isatty {
		tinfo, err := terminfo.LoadFromEnv()
		if err == nil {
			logger.tinfo = tinfo
		}
	}

	return logger
}

func (this *consoleLogger) Destroy() {
	/* -- nothing to do */
}

func (this *consoleLogger) LogObject(
	system string,
	subsystem Subsystem,
	severity Severity,
	verbosity Verbosity,
	object interface{},
) {
	/* -- the logger supports only line objects */
	line, ok := object.(LineObject)
	if !ok {
		return
	}

	/* -- lock the mutex avoiding inter-mixing of lines from different
	   threads. */
	this.mutex.Lock()
	defer this.mutex.Unlock()

	/* -- change the text color */
	if this.tinfo != nil {
		/* -- select color */
		var color int
		switch severity {
		case Critical:
			color = 1 /* -- red */
		case Error:
			color = 3 /* -- yellow */
		case Warning:
			color = 2 /* -- green */
		case Info:
			color = -1
		case Debug:
			color = -1
		default:
			color = -1
		}

		/* -- colorize the output */
		if color > 0 {
			this.output.WriteString(
				this.tinfo.Printf(terminfo.SetAForeground, color))
		}
	}

	/* -- print the message */
	fmt.Fprintf(
		this.output,
		"[%8s, %d] (%s): %s\n",
		severity.Code(),
		verbosity,
		subsystem,
		line.GetLogLine())

	/* -- reset terminal attributes back again */
	if this.tinfo != nil {
		this.output.WriteString(
			this.tinfo.Printf(terminfo.ExitAttributeMode))
	}
}
