// Package log
//
// @author: xwc1125
// @date: 2021/7/12
package log

import "os"

// InitLogs init log
func InitLogs(logConfig *LogConfig) {
	if logConfig.File.Save {
		handler := LvlFilterHandler(
			Lvl(logConfig.File.Level),
			Must.RotatingFileHandler(logConfig.File.GetLogFile(), 102400, JSONFormat()),
		)
		InitLogsWithFormat(logConfig, handler)
	} else {
		InitLogsWithFormat(logConfig, nil)
	}
}

// InitLogsWithFormat init log
func InitLogsWithFormat(logConfig *LogConfig, handler ...Handler) {
	// console setting
	stream := StreamHandler(os.Stderr, TerminalFormat(logConfig.Console.UseColor))
	consoleLogger := NewGlogHandler(stream)
	consoleLogger.Verbosity(Lvl(logConfig.Console.Level))
	consoleLogger.VModules(logConfig.Console.GetModules())

	PrintOrigins(logConfig.Console.ShowPath)

	// file setting
	if logConfig.File.Save {
		handlers := make([]Handler, 0)
		handlers = append(handlers, handler...)
		if logConfig.Console.Console {
			//StreamHandler(os.Stderr, LogfmtFormat()),
			handlers = append(handlers, consoleLogger)
		}
		Root().SetHandler(MultiHandler(
			handlers...,
		))
		Info("logs path", "path", logConfig.File.GetLogFile())
	} else {
		if logConfig.Console.Console {
			Root().SetHandler(consoleLogger)
		}
	}
}
