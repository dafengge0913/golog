// +build !windows,!linux

package golog

func NewPrinter() *PlainPrinter {
	return NewPlainPrinter()
}
