package golog

import (
	"fmt"
)

const (
	ansiColorFormat = "\033[1;%dm%s\033[m"

	ansiColorRed    = 31
	ansiColorGreen  = 32
	ansiColorYellow = 33
	ansiColorBlue   = 34
	ansiColorWhite  = 37
)

func ansiPrint(level LevelType, str string) error {
	_, err := fmt.Print(fmt.Sprintf(ansiColorFormat, getColorByLevel(level), str))
	return err
}

func getColorByLevel(level LevelType) int {
	switch level {
	case LevelDebug:
		return ansiColorWhite
	case LevelInfo:
		return ansiColorGreen
	case LevelWarn:
		return ansiColorYellow
	case LevelError:
		return ansiColorRed
	default:
		return ansiColorWhite
	}
}
