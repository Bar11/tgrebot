// description: log
//
// @author: xwc1125
// @date: 2019/9/18
package tests

import (
	"fmt"
	"github.com/chain5j/log15"
	"os"
	"testing"
)

var (
	ostream log.Handler
	glogger *log.GlogHandler
)

func TestLog(t *testing.T) {
	//usecolor := (isatty.IsTerminal(os.Stderr.Fd()) || isatty.IsCygwinTerminal(os.Stderr.Fd())) && os.Getenv("TERM") != "dumb"
	//output := io.Writer(os.Stderr)
	//if usecolor {
	//	output = colorable.NewColorableStderr()
	//}
	//ostream = log.StreamHandler(output, log.TerminalFormat(usecolor))
	ostream = log.StreamHandler(os.Stderr, log.TerminalFormat(true))
	glogger = log.NewGlogHandler(ostream)
	glogger.Verbosity(log.Lvl(5))

	//glogger.Vmodule(*vmodule)
	log.Root().SetHandler(glogger)

	// all loggers can have key/value context
	srvlog := log.New("app/server")

	// all log messages can have key/value context
	//srvlog.Warn("abnormal conn rate", "rate", "curRate", "low", "lowRate", "high", "highRate")

	//srvlog.Print("aaaaaa")
	//srvlog.Print(111, 222, 333, "12345")
	log.PrintOrigins(false)
	testLog := log.New("send next round change", "aa", "你好")
	testLog.Debug("send", "sta", "aaa")
	testLog.Info("send next round change222")
	testLog.Error("send next round change222")

	testLog2 := log.New("pbftcore", "aa", "你好")
	testLog2.Debug("send next round change send next round change send next round change send next round change", "sta", "aaa")
	testLog2.Info("send ")
	testLog2.Error("send next round change222")

	// child loggers with inherited context
	connlog := srvlog.New("srvlog", "raddr", "127.0.0.1")
	connlog.Info("connection open")

	// lazy evaluation
	connlog.Debug("ping remote", "latency", log.Lazy{"pingRemote"})

	// flexible configuration
	srvlog.SetHandler(log.MultiHandler(
		log.StreamHandler(os.Stderr, log.LogfmtFormat()),
		log.LvlFilterHandler(
			log.LvlError,
			log.Must.FileHandler("errors.json", log.JSONFormat()))))
}

func TestFileLog(t *testing.T) {
	stream := log.StreamHandler(os.Stderr, log.TerminalFormat(true))
	gLogger := log.NewGlogHandler(stream)
	gLogger.Verbosity(log.Lvl(4))
	gLogger.VModules([]string{"*"})
	log.PrintOrigins(true) // whether print path

	//consoleStream := log.StreamHandler(os.Stderr, log.TerminalFormat(true))
	//consoleLogger := log.NewGlogHandler(consoleStream)
	//consoleLogger.Verbosity(log.Lvl(3))
	//consoleLogger.VModules([]string{"*"})
	log.Root().SetHandler(gLogger)

	// save file
	gLogger.SetHandler(log.MultiHandler(
		stream,
		log.LvlFilterHandler(
			log.LvlDebug,
			log.Must.RotatingFileHandler("./logs",
				20,
				log.JSONFormat(),
			),
		),
	))

	for i := 0; ; i++ {
		log.Debug("Debug" + fmt.Sprintf("%d", i))
		if i/9 == 0 {
			log.Info("Info" + fmt.Sprintf("%d", i))
		}
		if i/11 == 0 {
			log.Error("Error" + fmt.Sprintf("%d", i))
		}
	}
}

func TestWriteLog(t *testing.T) {
	stream := log.StreamHandler(os.Stderr, log.TerminalFormat(true))
	gLogger := log.NewGlogHandler(stream)
	gLogger.Verbosity(log.Lvl(4))
	gLogger.VModules([]string{"*"})
	log.PrintOrigins(true) // whether print path

	log.Root().SetHandler(gLogger)

	// save file
	gLogger.SetHandler(log.MultiHandler(
		stream,
		log.LvlFilterHandler(
			log.LvlDebug,
			log.RotatingDayFileHandler(&log.TimeWriter{
				Dir:        "./logs",
				FileName:   "errors.json",
				Compress:   true,
				ReserveDay: 1,
			},
				log.JSONFormat()),
		),
	))

	for i := 0; ; i++ {
		log.Debug("Debug" + fmt.Sprintf("%d", i))
		if i/9 == 0 {
			log.Info("Info" + fmt.Sprintf("%d", i))
		}
		if i/11 == 0 {
			log.Error("Error" + fmt.Sprintf("%d", i))
		}
	}
}

func TestInit(t *testing.T) {
	log.InitLogs(&log.LogConfig{
		Console: log.ConsoleLogConfig{
			Level:    5,
			Modules:  "*",
			ShowPath: true,
			UseColor: true,
			Console:  false,
		},
		File: log.FileLogConfig{
			Level:    4,
			Modules:  "*",
			Save:     false,
			FilePath: "./logs",
			FileName: "errors.json",
		},
	})
	for i := 0; i < 1000; i++ {
		log.Debug("Debug" + fmt.Sprintf("%d", i))
		if i%9 == 0 {
			log.Info("Info" + fmt.Sprintf("%d", i))
		}
		if i%11 == 0 {
			log.Error("Error" + fmt.Sprintf("%d", i))
		}
	}
}

func TestInit2(t *testing.T) {
	logConfig := &log.LogConfig{
		Console: log.ConsoleLogConfig{
			Level:    5,
			Modules:  "*",
			ShowPath: true,
			UseColor: true,
			Console:  true,
		},
		File: log.FileLogConfig{
			Level:    3,
			Modules:  "*",
			Save:     true,
			FilePath: "./logs",
			FileName: "errors.json",
		},
	}
	log.InitLogsWithFormat(logConfig, log.LvlFilterHandler(
		log.Lvl(logConfig.File.Level),
		log.RotatingDayFileHandler(&log.TimeWriter{
			Dir:        logConfig.File.FilePath,
			FileName:   logConfig.File.FileName,
			Compress:   true,
			ReserveDay: 1,
		}, log.JSONFormat()),
	))
	for i := 0; i < 1000; i++ {
		log.Debug("Debug" + fmt.Sprintf("%d", i))
		if i%9 == 0 {
			log.Info("Info" + fmt.Sprintf("%d", i))
		}
		if i%11 == 0 {
			log.Error("Error" + fmt.Sprintf("%d", i))
		}
	}
}
