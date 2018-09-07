package golog

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type LevelType uint8

const (
	LevelDebug LevelType = iota
	LevelInfo
	LevelWarn
	LevelError
)

const (
	levelDebugMsg = "Debug"
	levelInfoMsg  = "Info"
	levelWarnMsg  = "Warn"
	levelErrorMsg = "Error"
)

const (
	defaultDateFormat           = "2006-01-02"
	defaultDateTimeFormat       = "2006-01-02 15:04:05.000"
	defaultLogFormatPrefixPrint = "[%-5s] [%s] : %s%s \n"
	defaultFuncSkip             = 3
)

var (
	pool *sync.Pool
)

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			return &logEntity{}
		},
	}
}

type config struct {
	isPrint        bool
	printColor     bool
	printPath      bool
	dateFormat     string
	dateTimeFormat string
	funcSkip       int
}

func NewConfig() *config {
	return &config{
		isPrint:        true,
		printColor:     true,
		printPath:      true,
		dateFormat:     defaultDateFormat,
		dateTimeFormat: defaultDateTimeFormat,
		funcSkip:       defaultFuncSkip,
	}
}

func (cfg *config) SetDateFormat(dateFormat string) {
	cfg.dateFormat = dateFormat
}

func (cfg *config) SetDateTimeFormat(dateTimeFormat string) {
	cfg.dateTimeFormat = dateTimeFormat
}

func (cfg *config) SetFuncSkip(skip int) {
	cfg.funcSkip = skip
}

func (cfg *config) SetIsPrint(isPrint bool) {
	cfg.isPrint = isPrint
}

func (log *Logger) SetPrintColor(printColor bool) {
	if log.printColor == printColor {
		return
	}
	log.printColor = printColor
	if printColor {
		log.printer = NewPrinter()
	} else {
		log.printer = NewPlainPrinter()
	}
}

func (cfg *config) SetPrintPath(printPath bool) {
	cfg.printPath = printPath
}

type IPrinter interface {
	Print(level LevelType, str string) error
}

type ILogWriter interface {
	Write(*logEntity) error
	Close() error
}

type Logger struct {
	*config
	level   LevelType
	writer  []ILogWriter
	printer IPrinter
}

type logEntity struct {
	msg    string
	level  LevelType
	time   time.Time
	caller string
}

func NewLogger(level LevelType, writer ...ILogWriter) *Logger {
	logger := &Logger{
		level:   level,
		writer:  writer,
		printer: NewPrinter(),
		config:  NewConfig(),
	}

	return logger
}

func (log *Logger) needPrint(level LevelType) bool {
	return log.isPrint && level >= log.level
}

func (log *Logger) doLog(level LevelType, msg string, args ...interface{}) {
	if !(log.needPrint(level) || len(log.writer) > 0) {
		// not print not write
		return
	}
	fMsg := fmt.Sprintf(msg, args...)
	t := time.Now()
	caller := log.getFuncCaller()
	if log.needPrint(level) {
		str := fmt.Sprintf(defaultLogFormatPrefixPrint, getLevelFlagMsg(level), log.getDateTimeStr(t), caller, fMsg)
		if err := log.printer.Print(level, str); err != nil {
			fmt.Print(str)
		}
	}
	for _, w := range log.writer {
		le := pool.Get().(*logEntity)
		le.msg = fMsg
		le.level = level
		le.time = t
		le.caller = caller
		if err := w.Write(le); err != nil {
			fmt.Errorf("[Write Log Error] :%v", err)
		}
	}
}

func (log *Logger) Debug(msg string, args ...interface{}) {
	log.doLog(LevelDebug, msg, args...)
}

func (log *Logger) Info(msg string, args ...interface{}) {
	log.doLog(LevelInfo, msg, args...)
}

func (log *Logger) Warn(msg string, args ...interface{}) {
	log.doLog(LevelWarn, msg, args...)
}

func (log *Logger) Error(msg string, args ...interface{}) {
	log.doLog(LevelError, msg, args...)
}

func (log *Logger) Close() {
	for _, w := range log.writer {
		if err := w.Close(); err != nil {
			fmt.Errorf("logger close error:%v", err)
		}
	}
}

func (log *Logger) getDateStr(t time.Time) string {
	return t.Format(log.dateFormat)
}

func (log *Logger) getDateTimeStr(t time.Time) string {
	return t.Format(log.dateTimeFormat)
}

func getLevelFlagMsg(level LevelType) string {
	switch level {
	case LevelDebug:
		return levelDebugMsg
	case LevelInfo:
		return levelInfoMsg
	case LevelWarn:
		return levelWarnMsg
	case LevelError:
		return levelErrorMsg
	default:
		return ""
	}
}

func (log *Logger) getFuncCaller() string {
	if !log.printPath {
		return ""
	}
	_, file, line, _ := runtime.Caller(log.funcSkip)
	return file + ":" + strconv.Itoa(line) + " -> "
}
