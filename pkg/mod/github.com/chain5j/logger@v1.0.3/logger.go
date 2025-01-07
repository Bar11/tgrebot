// Package logger
//
// @author: xwc1125
package logger

import (
	"strings"
	"sync"
)

type Logger interface {
	Name() string
	New(module string, ctx ...interface{}) Logger

	Trace(msg string, ctx ...interface{})
	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
	Crit(msg string, ctx ...interface{})

	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

var (
	root        Logger
	logAdapters = make([]Logger, 0)
	lock        sync.RWMutex
)

// RegisterLog 需要先注册logger
func RegisterLog(logger Logger) Logger {
	if logger == nil {
		return logger
	}
	lock.Lock()
	defer lock.Unlock()
	logAdapters = append(logAdapters, logger)
	root = logger
	return logger
}

// SwitchLog 切换日志kit
func SwitchLog(logKitName string) {
	lock.RLock()
	defer lock.RUnlock()
	for _, logger := range logAdapters {
		if strings.EqualFold(logger.Name(), logKitName) {
			root = logger
		}
	}
}

func getRoot() Logger {
	if root == nil {
		root = initLog()
	}
	return root
}

// New 方法用到root，因此此方法不能够应用到init()中，除非先执行了RegisterLog
func New(module string, ctx ...interface{}) Logger {
	return getRoot().New(module, ctx...)
}

func Trace(msg string, ctx ...interface{}) {
	getRoot().Trace(msg, ctx...)
}

func Debug(msg string, ctx ...interface{}) {
	getRoot().Debug(msg, ctx...)
}

func Info(msg string, ctx ...interface{}) {
	getRoot().Info(msg, ctx...)
}

func Warn(msg string, ctx ...interface{}) {
	getRoot().Warn(msg, ctx...)
}

func Error(msg string, ctx ...interface{}) {
	getRoot().Error(msg, ctx...)
}

func Crit(msg string, ctx ...interface{}) {
	getRoot().Crit(msg, ctx...)
}

func Printf(format string, v ...interface{}) {
	getRoot().Printf(format, v...)
}

func Print(v ...interface{}) {
	getRoot().Print(v...)
}

func Println(v ...interface{}) {
	getRoot().Println(v...)
}

func Fatal(v ...interface{}) {
	getRoot().Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	getRoot().Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	getRoot().Fatalln(v...)
}

func Panic(v ...interface{}) {
	getRoot().Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	getRoot().Panicf(format, v...)
}

func Panicln(v ...interface{}) {
	getRoot().Panicln(v)
}
