// Package log15
//
// @author: xwc1125
// @date: 2021/8/30
package main

import (
	log15 "github.com/chain5j/log15"
	"sync"
	"time"
)

func main() {
	initLogs(&log15.LogConfig{
		Console: log15.ConsoleLogConfig{
			Level:    4,
			Modules:  "*",
			ShowPath: false,
			UseColor: true,
			Console:  true,
		},
		File: log15.FileLogConfig{
			Level:    4,
			Modules:  "*",
			Save:     true,
			FilePath: "./logs",
			FileName: "error.json",
		},
	})
	var wg sync.WaitGroup
	startTime := CurrentTime()
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			log15.Info("=========================", "i", i)
			for j := 0; j < 10000; j++ {
				log15.Debug("test1 debug", "i", i, "j", j)
				if i%9 == 0 {
					log15.Info("test2 info", "i", i, "j", j)
				}
				if i%13 == 0 {
					log15.Error("test2 info", "i", i, "j", j)
				}
			}
		}(i)
	}
	wg.Wait()
	log15.Info("总耗时", "elapsed", CurrentTime()-startTime)
}

// 初始化日志
func initLogs(logConfig *log15.LogConfig) {
	// InitLogs init log
	if logConfig.File.Save {
		handler := log15.LvlFilterHandler(
			log15.Lvl(logConfig.File.Level),
			log15.Must.RotatingFileHandler(logConfig.File.GetLogFile(), 1024000*100, log15.JSONFormat()),
			//log15.Must.RotatingDayFileHandler(logConfig.File.FilePath, logConfig.File.FileName, true, 1, log15.JSONFormat()),
		)
		log15.InitLogsWithFormat(logConfig, handler)
	} else {
		log15.InitLogsWithFormat(logConfig, nil)
	}
}

func CurrentTime() int64 {
	return time.Now().UnixNano() / 1e6
}
