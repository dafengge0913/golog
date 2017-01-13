package golog

import (
	"fmt"
	"testing"
	"time"
)

func TestLogFile(t *testing.T) {
	writeCfg := NewLogWriterConfig()
	writeCfg.SetSaveInterval(time.Second)
	writer, err := NewLogWriterFile(LEVEL_INFO, "E:/dfg", "demo_log", true, writeCfg)
	if err != nil {
		t.Error("create log error:", err)
		return
	}
	log := NewLogger(LEVEL_DEBUG, writer)
	log.Debug("test info111")
	log.Info("test info2:%d : %s", 123, "abc")
	time.Sleep(time.Second * 2)
	log.SetDateTimeFormat(time.RFC1123)
	log.Warn("test info3")
	writer.SetDateFormat("2006@01@02")
	log.Error("test info4")
	if err := log.Close(); err != nil {
		fmt.Println("log close error: ", err)
	}
	fmt.Println("done!")
}
