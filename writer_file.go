package golog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	CACHE_SIZE             = 1024
	SAVE_INTERVAL          = time.Second * 3
	LOG_FORMAT_PREFIX_FILE = "[%-5s] [%s] : %s -> %s \n"
)

type LogWriterFile struct {
	bus            chan *LogEntity
	tickChan       *time.Ticker
	level          LevelType
	fileUrl        string
	curFileUrl     string
	rotate         bool
	rotateFileDate time.Time
}

func NewLogWriterFile(level LevelType, path string, fileName string, rotate bool) (ILogWriter, error) {
	if !pathExists(path) {
		return nil, fmt.Errorf("path not exists: %s", path)
	}

	fileUrl := filepath.Join(path, fileName)
	writer := &LogWriterFile{
		level:      level,
		fileUrl:    fileUrl,
		curFileUrl: fileUrl + ".log",
		rotate:     rotate,
		bus:        make(chan *LogEntity, CACHE_SIZE),
		tickChan:   time.NewTicker(SAVE_INTERVAL),
	}
	if rotate {
		writer.refreshRotateFile(time.Now())
	}
	go writer.serve()
	return writer, nil
}

func (w *LogWriterFile) serve() {
	for {
		select {
		case <-w.tickChan.C:
			if err := w.writeFile(); err != nil {
				fmt.Println("LogWriterFile log error: ", err)
			}
		}
	}
}

func (w *LogWriterFile) Write(logEntity *LogEntity) error {
	if logEntity.level < w.level {
		return nil
	}
	select {
	case w.bus <- logEntity:
		return nil
	default:
		w.writeFile()
		return errors.New("LogWriterFile bus overflow")
	}

}

func (w *LogWriterFile) Close() error {
	w.tickChan.Stop()
	if err := w.writeFile(); err != nil {
		fmt.Println("LogWriterFile log error: ", err)
	}
	close(w.bus)
	return nil
}

func (w *LogWriterFile) writeFile() error {
	file, err := os.OpenFile(w.curFileUrl, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	for len(w.bus) > 0 {
		logEntity, ok := <-w.bus
		if !ok {
			return errors.New("LogWriterFile bus is closed")
		}
		if w.rotate && logEntity.time.After(w.rotateFileDate) {
			w.refreshRotateFile(logEntity.time)
			file, err = os.OpenFile(w.curFileUrl, os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				return err
			}
		}
		fMsg := fmt.Sprintf(LOG_FORMAT_PREFIX_FILE, getLevelFlagMsg(logEntity.level), getDateTimeStr(logEntity.time), logEntity.caller, logEntity.msg)
		_, err = file.WriteString(fMsg)
		pool.Put(logEntity)
	}

	return err
}

func (w *LogWriterFile) refreshRotateFile(t time.Time) {
	w.rotateFileDate = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Add(time.Hour*24 - 1)
	w.curFileUrl = w.fileUrl + "_" + getDateStr(t) + ".log"
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
