// +build linux

package golog

import (
	"fmt"
)

const (
	LINUX_COLOR_FORMAT = "\033[1;%dm%s\033[m"

	COLOR_GRAY   = 30
	COLOR_RED    = 31
	COLOR_GREEN  = 32
	COLOR_YELLOW = 33
	COLOR_BLUE   = 34
	COLOR_WHITE  = 37
)

type Printer struct {
}

func NewPrinter() *Printer {
	printer := &Printer{}
	return printer
}

func (printer *Printer) Print(level LevelType, str string) error {
	_, err := fmt.Println(fmt.Sprintf(LINUX_COLOR_FORMAT, getColorByLevel(level), str))
	return err
}

func getColorByLevel(level LevelType) int {
	switch level {
	case LEVEL_DEBUG:
		return COLOR_GRAY
	case LEVEL_INFO:
		return COLOR_GREEN
	case LEVEL_WARN:
		return COLOR_YELLOW
	case LEVEL_ERROR:
		return COLOR_RED
	default:
		return COLOR_WHITE
	}
}
