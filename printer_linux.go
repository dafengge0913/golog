// +build linux

package golog

type Printer struct {
}

func NewPrinter() *Printer {
	printer := &Printer{}
	return printer
}

func (printer *Printer) Print(level LevelType, str string) error {
	return ansiPrint(level, str)
}
