package golog

import (
	"fmt"
)

type PlainPrinter struct {
}

func NewPlainPrinter() *PlainPrinter {
	return &PlainPrinter{}
}

func (printer *PlainPrinter) Print(level LevelType, str string) error {
	_, err := fmt.Print(str)
	return err
}
