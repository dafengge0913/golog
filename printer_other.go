// +build !windows,!linux

package golog

import (
	"fmt"
)

type Printer struct {
}

func NewPrinter() *Printer {
	return &Printer{}
}

func (printer *Printer) Print(level LevelType, str string) error {
	_, err := fmt.Print(str)
	return err
}
