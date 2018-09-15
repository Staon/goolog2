package goolog2

import (
	"fmt"
	"time"
)

type lineFormatterDefault struct {
	short bool
}

// Create new default line formatter
//
// Parameters:
//     short: determine whether the short format should be used
func NewLineFormatterDefault(
	short bool,
) LineFormatter {
	return &lineFormatterDefault{
		short: short,
	}
}

func (this *lineFormatterDefault) FormatLine(
	writer FileWriter,
	now time.Time,
	system string,
	subsystem Subsystem,
	severity Severity,
	verbosity Verbosity,
	line string,
) {
	/* -- determine message color */
	var color Color
	switch severity {
	case Critical:
		color = RED
	case Error:
		color = YELLOW
	case Warning:
		color = BLUE
	}
	writer.ChangeColor(color)

	/* -- write the formatted message */
	if this.short {
		fmt.Fprintf(
			writer,
			"[%8s, %d] (%s): %s\n",
			severity.Code(),
			verbosity,
			subsystem,
			line)
	} else {
		fmt.Fprintf(
			writer,
			"%s %s [%8s, %d] (%s): %s\n",
			system,
			now.Format("2006-01-02T15:04:05"),
			severity.Code(),
			verbosity,
			subsystem,
			line)
	}

	/* -- reset the color back */
	writer.ResetColor()
}
