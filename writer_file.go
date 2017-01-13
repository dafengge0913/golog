package golog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultCacheSize           = 1024
	defaultSaveInterval        = time.Second * 3
	defaultLogFormatPrefixFile = "[%-5s] [%s] : %s -> %s \n"
)

type logWriterConfig struct {
	cacheSize      uint32
	saveInterval   time.Duration
	dateFormat     string
	dateTimeFormat string
}

type LogWriterFile struct {
	*logWriterConfig
	bus            chan *logEntity
	tickChan       *time.Ticker
	level          LevelType
	fileUrl        string
	curFileUrl     string
	rotate         bool
	rotateFileDate time.Time
}

func NewLogWriterConfig() *logWriterConfig {
	return &logWriterConfig{
		cacheSize:      defaultCacheSize,
		saveInterval:   defaultSaveInterval,
		dateFormat:     defaultDateFormat,
		dateTimeFormat: defaultDateTimeFormat,
	}
}

func (cfg *logWriterConfig) SetCacheSize(cacheSize uint32) {
	cfg.cacheSize = cacheSize
}

func (cfg *logWriterConfig) SetSaveInterval(saveInterval time.Duration) {
	cfg.saveInterval = saveInterval
}

func (cfg *logWriterConfig) SetDateFormat(dateFormat string) {
	cfg.dateFormat = dateFormat
}

func (cfg *logWriterConfig) SetDateTimeFormat(dateTimeFormat string) {
	cfg.dateTimeFormat = dateTimeFormat
}

// config - using default config when nil
func NewLogWriterFile(level LevelType, path string, fileName string, rotate bool, config *logWriterConfig) (*LogWriterFile, error) {
	if !pathExists(path) {
		return nil, fmt.Errorf("path not exists: %s", path)
	}
	if config == nil {
		config = NewLogWriterConfig()
	}
	fileUrl := filepath.Join(path, fileName)
	writer := &LogWriterFile{
		logWriterConfig: config,
		bus:             make(chan *logEntity, config.cacheSize),
		tickChan:        time.NewTicker(config.saveInterval),
		level:           level,
		fileUrl:         fileUrl,
		curFileUrl:      fileUrl + ".log",
		rotate:          rotate,
	}
	writer.refreshRotateDate(time.Now())
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

func (w *LogWriterFile) Write(logEntity *logEntity) error {
	if logEntity.level < w.level {
		return nil
	}
	select {
	case w.bus <- logEntity:
		return nil
	default:
		w.writeFile()
		w.Write(logEntity)
		return errors.New("LogWriterFile bus overflow")
	}

}

func (w *LogWriterFile) Close() error {
	w.tickChan.Stop()
	err := w.writeFile()
	if err != nil {
		fmt.Println("LogWriterFile log error: ", err)
	}
	close(w.bus)
	return err
}

func (w *LogWriterFile) writeFile() error {
	if len(w.bus) == 0 {
		return nil
	}
	file, err := os.OpenFile(w.curFileUrl, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	for len(w.bus) > 0 {
		select {
		case logEntity, ok := <-w.bus:
			if !ok {
				return errors.New("LogWriterFile bus is closed")
			}
			if w.rotate && logEntity.time.After(w.rotateFileDate) {
				w.refreshRotateDate(logEntity.time)
				file, err = os.OpenFile(w.curFileUrl, os.O_APPEND|os.O_CREATE, 0666)
				if err != nil {
					return err
				}
			}
			fMsg := fmt.Sprintf(defaultLogFormatPrefixFile, getLevelFlagMsg(logEntity.level), w.getDateTimeStr(logEntity.time), logEntity.caller, logEntity.msg)
			_, err = file.WriteString(fMsg)
			pool.Put(logEntity)
		default:
			return nil
		}

	}

	return err
}

func (w *LogWriterFile) refreshRotateDate(t time.Time) {
	if w.rotate {
		w.rotateFileDate = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()).Add(time.Hour*24 - 1)
		w.refreshRotateFile(t)
	}
}

func (w *LogWriterFile) refreshRotateFile(t time.Time) {
	if w.rotate {
		w.curFileUrl = w.fileUrl + "_" + w.getDateStr(t) + ".log"
	}
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func (w *LogWriterFile) getDateStr(t time.Time) string {
	return t.Format(w.dateFormat)
}

func (w *LogWriterFile) getDateTimeStr(t time.Time) string {
	return t.Format(w.dateTimeFormat)
}

func (w *LogWriterFile) SetDateFormat(dateFormat string) {
	w.logWriterConfig.SetDateFormat(dateFormat)
	w.refreshRotateFile(w.rotateFileDate)
}
