// Package zap
//
// @author: xwc1125
package zap

import (
	"fmt"
	"os"
	"time"

	"github.com/chain5j/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	_    logger.Logger = new(log)
	root               = &log{
		Logger: getZapLogger(logger.DefaultLogConfig),
		fields: make([]zap.Field, 0),
	}
)

func init() {
	logger.RegisterLog(root)
}

type log struct {
	*zap.Logger // 只能输出结构化日志，但是性能要高于 SugaredLogger
	// sugarLogger *zap.SugaredLogger // 可以输出 结构化日志、非结构化日志。性能低于 zap.Logger
	fields []zap.Field
}

// InitWithConfig 初始化日志 logger
func InitWithConfig(config *logger.LogConfig) logger.Logger {
	zapLogger := getZapLogger(config)
	*root.Logger = *zapLogger
	return root
}

func getZapLogger(config *logger.LogConfig) *zap.Logger {
	var callerKey string
	if config.Console.ShowPath {
		callerKey = logger.FieldKeyFile
	}
	zapConfig := zapcore.EncoderConfig{
		MessageKey: logger.FieldKeyMsg,   // 结构化（json）输出：msg的key
		LevelKey:   logger.FieldKeyLevel, // 结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:    logger.FieldKeyTime,  // 结构化（json）输出：时间的key（INFO，WARN，ERROR等）
		NameKey:    logger.FieldKeyModule,

		CallerKey: callerKey, // 结构化（json）输出：打印日志的文件对应的Key
		// EncodeLevel:  zapcore.CapitalColorLevelEncoder, //将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: zapcore.ShortCallerEncoder, // 采用短文件路径编码输出（test/main.go:14 ）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("[" + t.Format(logger.DefaultTermTimestampFormat) + "]")
		}, // 输出的时间格式
		// EncodeName: func(loggerName string, enc zapcore.PrimitiveArrayEncoder) {
		//	enc.AppendString(loggerName)
		// },
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}

	filConfig := zapcore.EncoderConfig{
		MessageKey: logger.FieldKeyMsg,   // 结构化（json）输出：msg的key
		LevelKey:   logger.FieldKeyLevel, // 结构化（json）输出：日志级别的key（INFO，WARN，ERROR等）
		TimeKey:    logger.FieldKeyTime,  // 结构化（json）输出：时间的key（INFO，WARN，ERROR等）

		EncodeLevel:  zapcore.LowercaseLevelEncoder, // 将日志级别转换成大写（INFO，WARN，ERROR等）
		EncodeCaller: zapcore.ShortCallerEncoder,    // 采用短文件路径编码输出（test/main.go:14 ）
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(logger.DefaultTimestampFormat))
		}, // 输出的时间格式
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	}

	consoleLogLevel := getZapLevel(config.Console.Level)
	// 自定义日志级别：自定义Info级别
	fileLogLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= getZapLevel(config.File.Level)
	})

	// 获取io.Writer的实现
	fileWriter, _ := logger.GetWriter(config.File)
	// 实现多个输出
	tees := make([]zapcore.Core, 0)
	opts := make([]zap.Option, 0)
	if config.Console.Console {
		// NewConsoleEncoder 是非结构化输出
		tees = append(tees, zapcore.NewCore(NewConsoleEncoder(config, zapConfig), zapcore.AddSync(os.Stdout), consoleLogLevel))
	}
	if config.File.Save {
		// 同时将日志输出到控制台，NewJSONEncoder 是结构化输出
		tees = append(tees, zapcore.NewCore(zapcore.NewJSONEncoder(filConfig), zapcore.NewMultiWriteSyncer(zapcore.AddSync(fileWriter)), fileLogLevel))
	}
	if config.Console.ShowPath {
		opts = append(opts, zap.AddCaller())
	}
	core := zapcore.NewTee(tees...)

	zapLog := zap.New(core, opts...)
	return zapLog
}

func (l *log) Name() string {
	return "zap"
}

func (l *log) New(module string, ctx ...interface{}) logger.Logger {
	newLog := l.Logger.Named(module)
	fields := append(l.fields, l.getFields(ctx)...)
	newLog = newLog.With(fields...)
	return &log{
		Logger: newLog,
		// sugarLogger: newLog.Sugar(),
		fields: fields,
	}
}

func (l *log) Trace(msg string, ctx ...interface{}) {
	l.Logger.Debug(msg, l.getFields(ctx)...)
}

func (l *log) Debug(msg string, ctx ...interface{}) {
	l.Logger.Debug(msg, l.getFields(ctx)...)
}

func (l *log) getFields(ctx []interface{}) []zap.Field {
	if ctx != nil && len(ctx) > 0 {
		if len(ctx)%2 != 0 {
			ctx = append(ctx, "Null")
		}
	}

	files := make([]zap.Field, 0)
	for i := 0; i < len(ctx); i = i + 2 {
		files = append(files, zap.Any(logger.ToString(ctx[i]), ctx[i+1]))
	}
	return files
}

func (l *log) Info(msg string, ctx ...interface{}) {
	l.Logger.Info(msg, l.getFields(ctx)...)
}

func (l *log) Warn(msg string, ctx ...interface{}) {
	l.Logger.Warn(msg, l.getFields(ctx)...)
}

func (l *log) Error(msg string, ctx ...interface{}) {
	l.Logger.Error(msg, l.getFields(ctx)...)
}

func (l *log) Crit(msg string, ctx ...interface{}) {
	l.Logger.Fatal(msg, l.getFields(ctx)...)
}

func (l *log) Printf(format string, v ...interface{}) {
	l.Logger.Info(fmt.Sprintf(format, v...))
}

func (l *log) Print(v ...interface{}) {
	l.Logger.Info(fmt.Sprint(v...))
}

func (l *log) Println(v ...interface{}) {
	l.Logger.Info(fmt.Sprintln(v...))
}

func (l *log) Fatal(v ...interface{}) {
	l.Logger.Fatal(fmt.Sprint(v...))
}

func (l *log) Fatalf(format string, v ...interface{}) {
	l.Logger.Fatal(fmt.Sprintf(format, v...))
}

func (l *log) Fatalln(v ...interface{}) {
	l.Logger.Fatal(fmt.Sprintln(v...))
}

func (l *log) Panic(v ...interface{}) {
	l.Logger.Panic(fmt.Sprint(v...))
}

func (l *log) Panicf(format string, v ...interface{}) {
	l.Logger.Panic(fmt.Sprintf(format, v...))
}

func (l *log) Panicln(v ...interface{}) {
	l.Logger.Panic(fmt.Sprintln(v...))
}
