// +build linux

package golog

import (
	"fmt"
)

type Printer struct {
}

func NewPrinter() *Printer {
	printer := &Printer{}
	return printer
}

func (printer *Printer) Print(level LevelType, str string) error {
	return ansiPrint(level, str)
}
