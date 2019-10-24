package golog

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestLogFile(t *testing.T) {
	writeCfg := NewLogWriterConfig()
	writeCfg.SetSaveInterval(time.Second)
	writer, err := NewLogWriterFile(LevelInfo, "output", "demo_log", true, writeCfg)
	if err != nil {
		t.Error("create log error:", err)
		return
	}
	log := NewLogger(LevelDebug, writer)
	log.Debug("test info111")
	log.Info("test info2:%d : %s", 123, "abc")
	time.Sleep(time.Second * 2)
	log.SetDateTimeFormat(time.RFC1123)
	log.Warn("test info3")
	writer.SetDateFormat("2006@01@02")
	log.SetFuncSkip(2)
	log.Error("test info4")
	log.Close()
	fmt.Println("done!")
}

func TestLogFileConcurrent(t *testing.T) {
	writeCfg := NewLogWriterConfig()
	writeCfg.SetSaveInterval(time.Second)
	writer, err := NewLogWriterFile(LevelInfo, "E:/dfg", "demo_log", true, writeCfg)
	if err != nil {
		t.Error("create log error:", err)
		return
	}
	log := NewLogger(LevelDebug, writer)
	wg := &sync.WaitGroup{}
	for r := 0; r < 1000; r++ {
		wg.Add(1)
		go func(r int) {
			for i := 0; i < 100; i++ {
				log.Info("r:%d, i:%d", r, i)
			}
			wg.Done()
		}(r)
	}
	wg.Wait()
	log.Close()
	fmt.Println("done!")
}

func TestMultiLogWriter(t *testing.T) {
	writeCfg := NewLogWriterConfig()
	writeCfg.SetSaveInterval(time.Second)
	writer, err := NewLogWriterFile(LevelInfo, "output", "demo_log", true, writeCfg)

	errWriteCfg := NewLogWriterConfig()
	errWriteCfg.SetSaveInterval(time.Second)
	errWriter, err := NewLogWriterFile(LevelError, "output", "demo_log_error", true, errWriteCfg)
	if err != nil {
		t.Error("create log error:", err)
		return
	}
	log := NewLogger(LevelInfo, writer, errWriter)
	for i := 0; i < 1000; i++ {
		log.Info("test n=%d", i)
		log.Error("test error n=%d", i)
	}
	log.Close()
}

func TestColor(t *testing.T) {
	log := NewLogger(LevelDebug)
	log.Debug("Debug msg")
	log.Info("Info msg")
	log.Warn("Warn msg")
	log.Error("Error msg")
	log.SetPrintColor(false)
	log.Debug("Debug msg")
	log.Info("Info msg")
	log.Warn("Warn msg")
	log.Error("Error msg")
	log.SetPrintColor(true)
	log.Debug("Debug msg")
	log.Info("Info msg")
	log.Warn("Warn msg")
	log.Error("Error msg")
}

func TestPrintPath(t *testing.T) {
	log := NewLogger(LevelDebug)
	log.Info("Info msg")
	log.SetPrintPath(false)
	log.Info("Info msg")
}

var n = 0

func BenchmarkPrint(b *testing.B) {
	log := NewLogger(LevelInfo)
	for i := 0; i < b.N; i++ {
		n++
		log.Info("benchmark n=%d, n=%d, n=%d, n=%d, n=%d, n=%d, n=%d, n=%d", n, n, n, n, n, n, n, n)
	}
}

func BenchmarkWriterFile(b *testing.B) {
	writeCfg := NewLogWriterConfig()
	writeCfg.SetSaveInterval(time.Second)
	writer, err := NewLogWriterFile(LevelDebug, "output", "demo_log", true, writeCfg)
	if err != nil {
		b.Error("create log error:", err)
		return
	}
	log := NewLogger(LevelInfo, writer)
	log.SetIsPrint(false)
	for i := 0; i < b.N; i++ {
		n++
		log.Info("benchmark n=%d", n)
		if n%100 == 0 {
			log.Error("benchmark error n=%d", n)
		}
	}
	log.Close()
}
