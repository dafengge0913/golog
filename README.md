# golog
Simple logging library for Golang

## Installation  

`go get -u github.com/dafengge0913/golog`

## Quick start  

```golang
package main

import "golog"

func main() {
	log := golog.NewLogger(golog.LevelDebug)
	log.Debug("Debug msg")
	log.Info("Info msg, id=%d, name=%s", 1, "dfg")
	log.Warn("Warn msg")
	log.Error("Error msg")
}

```
### Output
![](https://i.loli.net/2018/09/07/5b9224e2ea924.png)

## Color  

Supported in Windows and Linux for now, This feature could be canceled by this setting :  
```golang
log.SetPrintColor(false)
```
and also could be activate anytime you want by set this to `true`.

## File path
File path and line number are optional, If you don't want this, do that :
```golang
log.SetPrintPath(false)
```
## Write to File

```golang
package main

import (
	"time"
	"golog"
)

func main() {
	writeCfg := golog.NewLogWriterConfig()
	// cache size of chan
	writeCfg.SetCacheSize(1024)
	// interval of flush to disk
	writeCfg.SetSaveInterval(time.Second * 5)
	// format date as a suffix to log file name when rotate
	writeCfg.SetDateFormat("2006@01@02")
	writer, err := golog.NewLogWriterFile(golog.LevelInfo, "D:/dfg", "demo_log", true, writeCfg)
	if err != nil {
		panic(err)
	}
	// only log error
	errWriter, err := golog.NewLogWriterFile(golog.LevelError, "D:/dfg", "demo_log_error", false, nil)
	if err != nil {
		panic(err)
	}
	log := golog.NewLogger(golog.LevelInfo, writer, errWriter)
	for i := 0; i < 100; i++ {
		log.Info("info i=%d", i)
		if i%10 == 0 {
			log.Error("error i=%d", i)
		}
	}
	// !! important
	log.Close()
}
```
You can create a file writer by call `golog.NewLogWriterFile`, For performance purpose, `LogWriterFile` is an asynchronous log recorder. when you spawn log message, it will be temporary stored in a chan, and flush to disk by a goroutine later.
**So it's important to call `Close()` before your application exit, make sure all message has been flushed in disk**

If you use all default settings for `LogWriterFile`, there is no need for creating a `logWriterConfig`, you can just pass a `nil` like this :
```golang
writer, err := NewLogWriterFile(LevelDebug, "D:/dfg", "demo_log", true, nil)
```

## License
Under the [MIT License](https://github.com/dafengge0913/golog/blob/master/LICENSE)