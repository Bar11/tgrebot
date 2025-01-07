// Package logger
//
// @author: xwc1125
package logger

import (
	log "github.com/chain5j/log15"
)

var (
	_ Logger = new(logDefault)
)

type logDefault struct {
	logger log.Logger
}

func initLog() Logger {
	log.InitLogsWithFormat(&log.LogConfig{
		Console: log.ConsoleLogConfig{
			Level:    int(DefaultLogConfig.Console.Level),
			Modules:  DefaultLogConfig.Console.Modules,
			ShowPath: DefaultLogConfig.Console.ShowPath,
			Format:   DefaultLogConfig.Console.Format,
			UseColor: DefaultLogConfig.Console.UseColor,
			Console:  DefaultLogConfig.Console.Console,
		},
	})
	return &logDefault{log.Root()}
}

func (l *logDefault) Name() string {
	return "logDefault"
}

func (l *logDefault) New(module string, ctx ...interface{}) Logger {
	return &logDefault{
		log.New(module, ctx...),
	}
}

func (l *logDefault) Trace(msg string, ctx ...interface{}) {
	if l.logger != nil {
		l.logger.Trace(msg, ctx...)
	} else {
		log.Trace(msg, ctx...)
	}
}

func (l *logDefault) Debug(msg string, ctx ...interface{}) {
	if l.logger != nil {
		l.logger.Debug(msg, ctx...)
	} else {
		log.Debug(msg, ctx...)
	}
}

func (l *logDefault) Info(msg string, ctx ...interface{}) {
	if l.logger != nil {
		l.logger.Info(msg, ctx...)
	} else {
		log.Info(msg, ctx...)
	}
}

func (l *logDefault) Warn(msg string, ctx ...interface{}) {
	if l.logger != nil {
		l.logger.Warn(msg, ctx...)
	} else {
		log.Warn(msg, ctx...)
	}
}

func (l *logDefault) Error(msg string, ctx ...interface{}) {
	if l.logger != nil {
		l.logger.Error(msg, ctx...)
	} else {
		log.Error(msg, ctx...)
	}
}

func (l *logDefault) Crit(msg string, ctx ...interface{}) {
	if l.logger != nil {
		l.logger.Error(msg, ctx...)
	} else {
		root.Error(msg, ctx...)
	}
}

func (l *logDefault) Printf(format string, v ...interface{}) {
	if l.logger != nil {
		l.logger.Printf(format, v...)
	}
}

func (l *logDefault) Print(v ...interface{}) {
	if l.logger != nil {
		l.logger.Print(v...)
	}
}

func (l *logDefault) Println(v ...interface{}) {
	if l.logger != nil {
		l.logger.Println(v...)
	}
}

func (l *logDefault) Fatal(v ...interface{}) {
	if l.logger != nil {
		l.logger.Fatal(v...)
	}
}

func (l *logDefault) Fatalf(format string, v ...interface{}) {
	if l.logger != nil {
		l.logger.Fatalf(format, v...)
	}
}

func (l *logDefault) Fatalln(v ...interface{}) {
	if l.logger != nil {
		l.logger.Fatalln(v...)
	}
}

func (l *logDefault) Panic(v ...interface{}) {
	if l.logger != nil {
		l.logger.Panic(v...)
	}
}

func (l *logDefault) Panicf(format string, v ...interface{}) {
	if l.logger != nil {
		l.logger.Panicf(format, v...)
	}
}

func (l *logDefault) Panicln(v ...interface{}) {
	if l.logger != nil {
		l.logger.Panicln(v...)
	}
}
