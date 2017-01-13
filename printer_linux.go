// +build linux

package golog

import (
	"fmt"
)

const (
	linuxColorFormat = "\033[1;%dm%s\033[m"

	colorGray   = 30
	colorRed    = 31
	colorGreen  = 32
	colorYellow = 33
	colorBlue   = 34
	colorWhite  = 37
)

type Printer struct {
}

func NewPrinter() *Printer {
	printer := &Printer{}
	return printer
}

func (printer *Printer) Print(level LevelType, str string) error {
	_, err := fmt.Println(fmt.Sprintf(linuxColorFormat, getColorByLevel(level), str))
	return err
}

func getColorByLevel(level LevelType) int {
	switch level {
	case LEVEL_DEBUG:
		return colorGray
	case LEVEL_INFO:
		return colorGreen
	case LEVEL_WARN:
		return colorYellow
	case LEVEL_ERROR:
		return colorRed
	default:
		return colorWhite
	}
}
