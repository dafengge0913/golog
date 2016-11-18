package golog

import (
	"fmt"
	"testing"
)

func TestLogFile(t *testing.T) {
	writer, err := NewLogWriterFile(LEVEL_INFO, "/data/dfg", "demo_log", true)
	if err != nil {
		t.Error("create log error:", err)
		return
	}

	log := NewLogger(LEVEL_DEBUG, writer)

	log.Debug("test info111")
	log.Info("test info2:%d : %s", 123, "abc")
	log.Warn("test info3")
	log.Error("test info4")
	if err := log.Close(); err != nil {
		fmt.Println("log close error: ", err)
	}
	fmt.Println("done!")
}
