// Package logger
//
// @author: xwc1125
package logger

import (
	"os"
	"path/filepath"
	"strings"
)

type Lvl int

const (
	LvlFatal Lvl = iota
	LvlError
	LvlWarn
	LvlInfo
	LvlDebug
	LvlTrace
)

func (l Lvl) String() string {
	switch l {
	case LvlTrace:
		return "trace"
	case LvlDebug:
		return "debug"
	case LvlInfo:
		return "info"
	case LvlWarn:
		return "warn"
	case LvlError:
		return "error"
	case LvlFatal:
		return "fatal"
	default:
		panic("bad level")
	}
}

var DefaultLogConfig = &LogConfig{
	Console: ConsoleLogConfig{
		Level:    3,
		Modules:  "*",
		ShowPath: false,
		UseColor: false,
		Console:  true,
	},
}

// LogConfig log config
type LogConfig struct {
	Console ConsoleLogConfig `json:"console" mapstructure:"console"`
	File    FileLogConfig    `json:"file" mapstructure:"file"`
}

type ConsoleLogConfig struct {
	Level    Lvl    `json:"level" mapstructure:"level"`         // log level
	Modules  string `json:"modules" mapstructure:"modules"`     // need to show modules。"*":all
	ShowPath bool   `json:"show_path" mapstructure:"show_path"` // show path
	Format   string `json:"format" mapstructure:"format"`       // format

	UseColor bool `json:"use_color" mapstructure:"use_color"` // show console color
	Console  bool `json:"console" mapstructure:"console"`     // show console
}

type FileLogConfig struct {
	Level  Lvl    `json:"level" mapstructure:"level"`   // log level
	Format string `json:"format" mapstructure:"format"` // format

	Save     bool   `json:"save" mapstructure:"save"`           // whether save file
	FilePath string `json:"file_path" mapstructure:"file_path"` // filepath
	FileName string `json:"file_name" mapstructure:"file_name"` // filename prefix

	MaxAge       int64 `json:"max_age" mapstructure:"max_age"`             // 文件最大保存时间[小时]
	RotationTime int64 `json:"rotation_time" mapstructure:"rotation_time"` // 日志切割时间间隔[小时]
}

func (c *ConsoleLogConfig) GetModules() []string {
	if c.Modules == "" {
		return []string{"*"}
	}
	return strings.Split(c.Modules, ",")
}

func (c *FileLogConfig) GetLogFile() string {
	if !c.Save {
		return ""
	}
	file := c.FileName
	if file == "" {
		file = "log.json"
	}
	if c.FilePath != "" {
		os.MkdirAll(c.FilePath, os.ModePerm)
		file = filepath.Join(c.FilePath, file)
	} else {
		file = filepath.Join(".", file)
	}
	return file
}

func (c *FileLogConfig) GetMaxAge() int64 {
	if c.MaxAge <= 0 {
		return 24 * 30 // 30天
	}
	return c.MaxAge
}

func (c *FileLogConfig) GetRotationTime() int64 {
	if c.RotationTime <= 0 {
		return 24 // 24小时
	}
	return c.RotationTime
}
