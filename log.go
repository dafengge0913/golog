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
	LEVEL_DEBUG LevelType = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
)

const (
	LEVEL_DEBUG_MSG = "Debug"
	LEVEL_INFO_MSG  = "Info"
	LEVEL_WARN_MSG  = "Warn"
	LEVEL_ERROR_MSG = "Error"
)

const (
	DATE_FORMAT             = "2006-01-02"
	DATE_TIME_FORMAT        = "2006-01-02 15:04:05.000"
	LOG_FORMAT_PREFIX_PRINT = "[%-5s] [%s] : %s -> %s \n"
)

var (
	pool *sync.Pool
)

func init() {
	pool = &sync.Pool{
		New: func() interface{} {
			return &LogEntity{}
		},
	}
}

type IPrinter interface {
	Print(level LevelType, str string) error
}

type ILogWriter interface {
	Write(*LogEntity) error
	Close() error
}

type Logger struct {
	level   LevelType
	writer  ILogWriter
	printer IPrinter
}

type LogEntity struct {
	msg    string
	level  LevelType
	time   time.Time
	caller string
}

func NewLogger(level LevelType, writer ILogWriter) *Logger {
	logger := &Logger{
		level:   level,
		writer:  writer,
		printer: NewPrinter(),
	}

	return logger
}

func (log *Logger) doLog(level LevelType, msg string, args ...interface{}) {
	fMsg := fmt.Sprintf(msg, args...)
	t := time.Now()
	caller := getFuncCaller(3)
	if level >= log.level {
		log.printer.Print(level, fmt.Sprintf(LOG_FORMAT_PREFIX_PRINT, getLevelFlagMsg(level), getDateTimeStr(t), caller, fMsg))
	}
	if log.writer != nil {
		le := pool.Get().(*LogEntity)
		le.msg = fMsg
		le.level = level
		le.time = t
		le.caller = caller
		if err := log.writer.Write(le); err != nil {
			fmt.Println("[Write Log Error] :", err)
		}
	}
}

func (log *Logger) Debug(msg string, args ...interface{}) {
	log.doLog(LEVEL_DEBUG, msg, args...)
}

func (log *Logger) Info(msg string, args ...interface{}) {
	log.doLog(LEVEL_INFO, msg, args...)
}

func (log *Logger) Warn(msg string, args ...interface{}) {
	log.doLog(LEVEL_WARN, msg, args...)
}

func (log *Logger) Error(msg string, args ...interface{}) {
	log.doLog(LEVEL_ERROR, msg, args...)
}

func (log *Logger) Close() error {
	if log.writer != nil {
		return log.writer.Close()
	}
	return nil
}

func getDateStr(t time.Time) string {
	return t.Format(DATE_FORMAT)
}

func getDateTimeStr(t time.Time) string {
	return t.Format(DATE_TIME_FORMAT)
}

func getLevelFlagMsg(level LevelType) string {
	switch level {
	case LEVEL_DEBUG:
		return LEVEL_DEBUG_MSG
	case LEVEL_INFO:
		return LEVEL_INFO_MSG
	case LEVEL_WARN:
		return LEVEL_WARN_MSG
	case LEVEL_ERROR:
		return LEVEL_ERROR_MSG
	default:
		return ""
	}
}

func getFuncCaller(n int) string {
	_, file, line, _ := runtime.Caller(n)
	return file + ":" + strconv.Itoa(line)
}
