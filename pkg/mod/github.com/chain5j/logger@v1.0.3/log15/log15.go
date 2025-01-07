// Package log15
//
// @author: xwc1125
package log15

import (
	log "github.com/chain5j/log15"
	"github.com/chain5j/logger"
)

var (
	_    logger.Logger = new(log15)
	root               = &log15{log.Root()}
)

func init() {
	setLog15(logger.DefaultLogConfig)
	logger.RegisterLog(root)
}

type log15 struct {
	logger log.Logger
}

func InitWithConfig(config *logger.LogConfig) logger.Logger {
	setLog15(config)

	return root
}

func setLog15(config *logger.LogConfig) {
	logConfig := &log.LogConfig{
		Console: log.ConsoleLogConfig{
			Level:    int(config.Console.Level),
			Modules:  config.Console.Modules,
			ShowPath: config.Console.ShowPath,
			Format:   config.Console.Format,
			UseColor: config.Console.UseColor,
			Console:  config.Console.Console,
		},
		File: log.FileLogConfig{
			Level:    int(config.File.Level),
			Format:   config.File.Format,
			Save:     config.File.Save,
			FilePath: config.File.FilePath,
			FileName: config.File.FileName,
		},
	}
	if config.File.Save {
		writer, _ := logger.GetWriter(config.File)
		handler := log.LvlFilterHandler(
			log.Lvl(config.File.Level),
			log.WriterHandler(writer, log.JSONFormat()),
		)
		log.InitLogsWithFormat(logConfig, handler)
	} else {
		log.InitLogsWithFormat(logConfig, nil)
	}
}

func (l *log15) Name() string {
	return "log15"
}

func (l *log15) New(module string, ctx ...interface{}) logger.Logger {
	return &log15{
		root.logger.New(module, ctx...),
	}
}

func (l *log15) Trace(msg string, ctx ...interface{}) {
	if l.logger != nil {
		l.logger.Trace(msg, ctx...)
	} else {
		log.Trace(msg, ctx...)
	}
}

func (l *log15) Debug(msg string, ctx ...interface{}) {
	if l.logger != nil {
		l.logger.Debug(msg, ctx...)
	} else {
		log.Debug(msg, ctx...)
	}
}

func (l *log15) Info(msg string, ctx ...interface{}) {
	if l.logger != nil {
		l.logger.Info(msg, ctx...)
	} else {
		log.Info(msg, ctx...)
	}
}

func (l *log15) Warn(msg string, ctx ...interface{}) {
	if l.logger != nil {
		l.logger.Warn(msg, ctx...)
	} else {
		log.Warn(msg, ctx...)
	}
}

func (l *log15) Error(msg string, ctx ...interface{}) {
	if l.logger != nil {
		l.logger.Error(msg, ctx...)
	} else {
		log.Error(msg, ctx...)
	}
}

func (l *log15) Crit(msg string, ctx ...interface{}) {
	if l.logger != nil {
		l.logger.Error(msg, ctx...)
	} else {
		root.Error(msg, ctx...)
	}
}

func (l *log15) Printf(format string, v ...interface{}) {
	if l.logger != nil {
		l.logger.Printf(format, v...)
	}
}

func (l *log15) Print(v ...interface{}) {
	if l.logger != nil {
		l.logger.Print(v...)
	}
}

func (l *log15) Println(v ...interface{}) {
	if l.logger != nil {
		l.logger.Println(v...)
	}
}

func (l *log15) Fatal(v ...interface{}) {
	if l.logger != nil {
		l.logger.Fatal(v...)
	}
}

func (l *log15) Fatalf(format string, v ...interface{}) {
	if l.logger != nil {
		l.logger.Fatalf(format, v...)
	}
}

func (l *log15) Fatalln(v ...interface{}) {
	if l.logger != nil {
		l.logger.Fatalln(v...)
	}
}

func (l *log15) Panic(v ...interface{}) {
	if l.logger != nil {
		l.logger.Panic(v...)
	}
}

func (l *log15) Panicf(format string, v ...interface{}) {
	if l.logger != nil {
		l.logger.Panicf(format, v...)
	}
}

func (l *log15) Panicln(v ...interface{}) {
	if l.logger != nil {
		l.logger.Panicln(v...)
	}
}
