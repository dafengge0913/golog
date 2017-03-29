// +build windows

package golog

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

const (
	colorWhite  uintptr = 7
	colorGray   uintptr = 8
	colorGreen  uintptr = 10
	colorBlue   uintptr = 11
	colorRed    uintptr = 12
	colorYellow uintptr = 14
)

type Printer struct {
	isTerminal                 bool
	kernel32                   *syscall.LazyDLL
	setConsoleTextAttribute    *syscall.LazyProc
	getConsoleScreenBufferInfo *syscall.LazyProc
}

func NewPrinter() *Printer {
	printer := &Printer{
		kernel32: syscall.NewLazyDLL("kernel32.dll"),
	}
	printer.isTerminal = printer.isWinTerminal()
	printer.setConsoleTextAttribute = printer.kernel32.NewProc("SetConsoleTextAttribute")
	printer.getConsoleScreenBufferInfo = printer.kernel32.NewProc("GetConsoleScreenBufferInfo")
	return printer
}

type coord struct {
	X, Y int16
}

type consoleScreenBufferInfo struct {
	dwSize              coord
	dwCursorPosition    coord
	wAttributes         uint16
	srWindow            smallRect
	dwMaximumWindowSize coord
}

type smallRect struct {
	Left, Top, Right, Bottom int16
}

func (printer *Printer) Print(level LevelType, str string) error {

	if !printer.isTerminal {
		fmt.Print(str)
		return nil
	}

	origin, err := printer.getColor()
	if err != nil {
		return err
	}

	if err := printer.setColor(printer.getColorByLevel(level)); err != nil {
		return err
	}

	fmt.Print(str)

	if err := printer.setColor(origin); err != nil {
		return err
	}
	return nil
}

func (printer *Printer) getColorByLevel(level LevelType) uintptr {
	switch level {
	case LevelDebug:
		return colorGray
	case LevelInfo:
		return colorGreen
	case LevelWarn:
		return colorYellow
	case LevelError:
		return colorRed
	default:
		return colorWhite
	}
}

func (printer *Printer) getColor() (uintptr, error) {
	var info consoleScreenBufferInfo
	r1, _, err := printer.getConsoleScreenBufferInfo.Call(uintptr(syscall.Stdout), uintptr(unsafe.Pointer(&info)))
	if int(r1) == 0 {
		return 0, err
	}
	return uintptr(info.wAttributes), nil
}

func (printer *Printer) setColor(color uintptr) error {
	r1, _, err := printer.setConsoleTextAttribute.Call(uintptr(syscall.Stdout), color)
	if int(r1) == 0 {
		return err
	}
	return nil
}

func (printer *Printer) isWinTerminal() bool {
	var getConsoleMode = printer.kernel32.NewProc("GetConsoleMode")
	var st uint32
	r, _, e := syscall.Syscall(getConsoleMode.Addr(), 2, os.Stdout.Fd(), uintptr(unsafe.Pointer(&st)), 0)
	return r != 0 && e == 0
}
