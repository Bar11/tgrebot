// Package logrus
//
// @author: xwc1125
// @date: 2021/8/30
package logrus

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/chain5j/logger"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	_    logger.Logger = new(log)
	root               = &log{Logger: logrus.New()}
)

func init() {
	InitWithConfig(logger.DefaultLogConfig)
	logger.RegisterLog(root)
}

type log struct {
	*logrus.Logger
	entry *logrus.Entry
}

func InitWithConfig(config *logger.LogConfig) logger.Logger {
	root.SetReportCaller(config.Console.ShowPath)
	root.SetLevel(getLevel(config.Console.Level))

	root.SetFormatter(&Formatter{
		UseColor:        config.Console.UseColor,
		LocationEnabled: config.Console.ShowPath,
		Modules:         strings.Split(config.Console.Modules, ","),
	})

	if config.File.Save {
		writer, _ := logger.GetWriter(config.File)
		var writerMap = make(map[logrus.Level]io.Writer, 0)
		if config.File.Level >= logger.LvlTrace {
			writerMap[logrus.TraceLevel] = writer
		}
		if config.File.Level >= logger.LvlDebug {
			writerMap[logrus.DebugLevel] = writer
		}
		if config.File.Level >= logger.LvlInfo {
			writerMap[logrus.InfoLevel] = writer
		}
		if config.File.Level >= logger.LvlWarn {
			writerMap[logrus.WarnLevel] = writer
		}
		if config.File.Level >= logger.LvlError {
			writerMap[logrus.ErrorLevel] = writer
		}
		if config.File.Level >= logger.LvlFatal {
			writerMap[logrus.FatalLevel] = writer
			writerMap[logrus.PanicLevel] = writer
		}
		root.Hooks.Add(lfshook.NewHook(
			lfshook.WriterMap(writerMap),
			&logrus.JSONFormatter{
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyTime:        logger.FieldKeyTime,
					logrus.FieldKeyLevel:       logger.FieldKeyLevel,
					logrus.FieldKeyMsg:         logger.FieldKeyMsg,
					logrus.FieldKeyFile:        logger.FieldKeyFile,
					logrus.FieldKeyLogrusError: logger.FieldKeyError,
					"module":                   logger.FieldKeyModule,
				},
			},
		))
	}
	if !config.Console.Console {
		root.SetOutput(ioutil.Discard)
	}

	return root
}

func (l *log) Name() string {
	return "logrus"
}

func (l *log) New(module string, ctx ...interface{}) logger.Logger {
	return New(module, ctx...)
}

func New(module string, ctx ...interface{}) logger.Logger {
	fields := getFields(ctx...)
	fields["module"] = module
	return &log{
		Logger: root.Logger,
		entry:  root.Logger.WithFields(fields),
	}
}

func getFields(ctx ...interface{}) map[string]interface{} {
	if ctx != nil && len(ctx) > 0 {
		if len(ctx)%2 != 0 {
			ctx = append(ctx, "NULL")
		}
	}
	var fields = make(map[string]interface{}, 0)
	for i := 0; i < len(ctx); i = i + 2 {
		fields[logger.ToString(ctx[i])] = ctx[i+1]
	}
	return fields
}

func (l *log) Trace(msg string, ctx ...interface{}) {
	if l.entry != nil {
		l.entry.WithFields(getFields(ctx...)).Trace(msg)
	} else {
		l.Logger.WithFields(getFields(ctx...)).Trace(msg)
	}
}

func (l *log) Debug(msg string, ctx ...interface{}) {
	if l.entry != nil {
		l.entry.WithFields(getFields(ctx...)).Debug(msg)
	} else {
		l.Logger.WithFields(getFields(ctx...)).Debug(msg)
	}
}

func (l *log) Info(msg string, ctx ...interface{}) {
	if l.entry != nil {
		l.entry.WithFields(getFields(ctx...)).Info(msg)
	} else {
		l.Logger.WithFields(getFields(ctx...)).Info(msg)
	}
}

func (l *log) Warn(msg string, ctx ...interface{}) {
	if l.entry != nil {
		l.entry.WithFields(getFields(ctx...)).Warn(msg)
	} else {
		l.Logger.WithFields(getFields(ctx...)).Warn(msg)
	}
}

func (l *log) Error(msg string, ctx ...interface{}) {
	if l.entry != nil {
		l.entry.WithFields(getFields(ctx...)).Error(msg)
	} else {
		l.Logger.WithFields(getFields(ctx...)).Error(msg)
	}
}

func (l *log) Crit(msg string, ctx ...interface{}) {
	if l.entry != nil {
		l.entry.WithFields(getFields(ctx...)).Fatal(msg)
	} else {
		l.Logger.WithFields(getFields(ctx...)).Fatal(msg)
	}
}
