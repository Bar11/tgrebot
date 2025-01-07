// Package test
//
// @author: xwc1125
package test

import (
	"log"
	"testing"

	"github.com/chain5j/logger"
	"github.com/chain5j/logger/log15"
	"github.com/chain5j/logger/logrus"
	"github.com/chain5j/logger/zap"
)

func TestZapLog(t *testing.T) {
	zapLog()
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>", "zapLog", "zapLog")
}
func TestLog15Log(t *testing.T) {
	log15Log()
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>", "log15Log", "log15Log")
}
func TestLogrusLog(t *testing.T) {
	logrusLog()
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>", "logrusLog", "logrusLog")
}

func log15Log() {
	log15.InitWithConfig(&logger.LogConfig{
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
			FileName: "logger.json",
		},
	})

	startTime := logger.CurrentTime()

	log := logger.New("log15")

	for i := 0; i < 2; i++ {
		log.Info("=========================", "i", i)
		for j := 0; j < 10; j++ {
			log.Debug("test1 debug", "i", i, "j", j)
			if j%3 == 0 {
				log.Info("test2 info", "i", i, "j", j)
			}
			if j%5 == 0 {
				logger.Error("test2 info", "i", i, "j", j)
			}
		}
	}
	logger.Info("总耗时", "elapsed", logger.CurrentTime()-startTime)
}

func zapLog() {
	zap.InitWithConfig(&logger.LogConfig{
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
			FileName: "logger.json",
		},
	})

	startTime := logger.CurrentTime()

	log := logger.New("zap")

	for i := 0; i < 2; i++ {
		log.Info("=========================", "i", i)
		for j := 0; j < 10; j++ {
			log.Debug("test1 debug", "i", i, "j", j)
			if j%3 == 0 {
				log.Info("test2 info", "i", i, "j", j)
			}
			if j%5 == 0 {
				logger.Error("test2 info", "i", i, "j", j)
			}
		}
	}
	logger.Info("总耗时", "elapsed", logger.CurrentTime()-startTime)
}

func logrusLog() {
	logrus.InitWithConfig(&logger.LogConfig{
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
			FileName: "logger.json",
		},
	})

	startTime := logger.CurrentTime()

	log := logger.New("logrus")

	for i := 0; i < 2; i++ {
		log.Info("=========================", "i", i)
		for j := 0; j < 10; j++ {
			log.Debug("test1 debug", "i", i, "j", j)
			if j%3 == 0 {
				log.Info("test2 info", "i", i, "j", j)
			}
			if j%5 == 0 {
				logger.Error("test2 info", "i", i, "j", j)
			}
		}
	}
	logger.Info("总耗时", "elapsed", logger.CurrentTime()-startTime)
}
