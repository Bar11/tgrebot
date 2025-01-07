// Package zap
//
// @author: xwc1125
package zap

import (
	"github.com/chain5j/logger"
	"go.uber.org/zap/zapcore"
)

func getZapLevel(lvl logger.Lvl) zapcore.Level {
	switch lvl {
	case logger.LvlTrace, logger.LvlDebug:
		return zapcore.DebugLevel
	case logger.LvlInfo:
		return zapcore.InfoLevel
	case logger.LvlWarn:
		return zapcore.WarnLevel
	case logger.LvlError:
		return zapcore.ErrorLevel
	case logger.LvlFatal:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
