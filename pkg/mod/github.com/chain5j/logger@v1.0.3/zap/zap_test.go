// Package zap
//
// @author: xwc1125
package zap

import (
	syslog "log"
	"sync"
	"testing"

	"github.com/chain5j/logger"
)

func TestLog(t *testing.T) {
	logConfig := &logger.LogConfig{
		Console: logger.ConsoleLogConfig{
			Level:    4,
			Modules:  "*",
			ShowPath: false,
			UseColor: true,
			Console:  true,
		},
		File: logger.FileLogConfig{
			Level:    4,
			Save:     true,
			FilePath: "./logs",
			FileName: "log.json",
		},
	}
	log1 := InitWithConfig(logConfig)

	var wg sync.WaitGroup
	startTime := logger.CurrentTime()

	debugLog := log1.New("Debug", "Debug", "Debug")
	infoLog := log1.New("Info", "Info", "Info")
	errorLog := log1.New("Error", "Error", "Error")

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			syslog.Println("=========================", "i", i)
			for j := 0; j < 100; j++ {
				debugLog.Debug("test1 debug", "i", i, "j", j)
				if i%9 == 0 {
					infoLog.Info("test2 info", "i", i, "j", j)
				}
				if i%13 == 0 {
					errorLog.Error("test2 info", "i", i, "j", j)
				}
			}
		}(i)
	}
	wg.Wait()
	log1.Info("总耗时", "elapsed", logger.CurrentTime()-startTime)
}
