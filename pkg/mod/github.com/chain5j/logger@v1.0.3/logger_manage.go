// Package logger
//
// @author: xwc1125
package logger

import (
	"reflect"
	"sync"
)

var (
	logsMap sync.Map // module==>logger.Logger
)

func Log(obj interface{}) Logger {
	module := getClsName(obj)
	if val, ok := logsMap.Load(module); ok {
		return val.(Logger)
	}
	l := New(module)
	logsMap.Store(module, l)
	return l
}

func getClsName(obj interface{}) string {
	switch obj.(type) {
	case struct{}:
		return reflect.TypeOf(obj).String()
	default:
		return ToString(obj)
	}
}
