// Package logrus
//
// @author: xwc1125
// @date: 2021/9/2
package logrus

import (
	"github.com/chain5j/logger"
	"github.com/sirupsen/logrus"
)

func getLevel(lvl logger.Lvl) logrus.Level {
	switch lvl {
	case logger.LvlTrace:
		return logrus.TraceLevel
	case logger.LvlDebug:
		return logrus.DebugLevel
	case logger.LvlInfo:
		return logrus.InfoLevel
	case logger.LvlWarn:
		return logrus.WarnLevel
	case logger.LvlError:
		return logrus.ErrorLevel
	case logger.LvlFatal:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}
