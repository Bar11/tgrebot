// Copyright (c) 2016 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package zap

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf8"

	"github.com/chain5j/logger"
	"github.com/chain5j/logger/zap/bufferpool"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

var _sliceEncoderPool = sync.Pool{
	New: func() interface{} {
		return &sliceArrayEncoder{elems: make([]interface{}, 0, 2)}
	},
}

func getSliceEncoder() *sliceArrayEncoder {
	return _sliceEncoderPool.Get().(*sliceArrayEncoder)
}

func putSliceEncoder(e *sliceArrayEncoder) {
	e.elems = e.elems[:0]
	_sliceEncoderPool.Put(e)
}

type consoleEncoder struct {
	config *logger.LogConfig
	*jsonEncoder

	modules    []string
	modulesReg atomic.Value
}

// NewConsoleEncoder creates an encoder whose output is designed for human -
// rather than machine - consumption. It serializes the core log entry data
// (message, level, timestamp, etc.) in a plain-text format and leaves the
// structured context as JSON.
//
// Note that although the console encoder doesn't use the keys specified in the
// encoder configuration, it will omit any element whose key is set to the empty
// string.
func NewConsoleEncoder(config *logger.LogConfig, cfg zapcore.EncoderConfig) zapcore.Encoder {
	return &consoleEncoder{
		config:      config,
		jsonEncoder: newJSONEncoder(cfg, false),
		modules:     config.Console.GetModules(),
	}
}

func (c *consoleEncoder) Clone() zapcore.Encoder {
	return &consoleEncoder{
		config:      c.config,
		jsonEncoder: c.jsonEncoder.Clone().(*jsonEncoder),
		modules:     c.modules,
		modulesReg:  c.modulesReg,
	}
}

func (c *consoleEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	line := bufferpool.Get()
	// line.AppendString(c.jsonEncoder.buf.String())
	module := ent.LoggerName
	if !c.isPrint(module) {
		return line, nil
	}

	arr := getSliceEncoder()
	levelColor := 0
	if c.config.Console.UseColor {
		levelColor = getColorByLevel(ent.Level)
	}
	// level
	var level = strings.ToUpper(ent.Level.String())
	{
		levelLen := len(level)
		if levelColor > 0 {
			level = logger.ColorRender(level, levelColor, 0)
		}
		line.AppendString(level)
		line.AppendString(strings.Repeat(" ", 5-levelLen))
	}

	// time
	{
		timestampFormat := c.config.Console.Format
		if timestampFormat == "" {
			timestampFormat = logger.DefaultTermTimestampFormat
		}

		line.AppendString(" ")
		line.AppendString("[")
		line.AppendString(ent.Time.Format(timestampFormat))
		line.AppendString("]")
	}
	// module
	paddingLen := 0
	{
		module := ent.LoggerName
		if levelColor > 0 {
			module = logger.ColorRender(module, levelColor, 0)
		}
		zapcore.FullNameEncoder(" ["+module+"]", arr)
		paddingLen = utf8.RuneCountInString(module)
	}

	// writer path
	{
		if ent.Caller.Defined && c.CallerKey != "" && c.EncodeCaller != nil {
			c.EncodeCaller(ent.Caller, arr)
		}
		for i := range arr.elems {
			if i > 0 {
				line.AppendByte('\t')
			}
			fmt.Fprint(line, arr.elems[i])
		}
		putSliceEncoder(arr)
	}

	// msg
	if c.MessageKey != "" {
		c.addTabIfNecessary(line)
		line.AppendString(ent.Message)
	}

	// write padding
	{
		// try to justify the log output for short messages
		length := utf8.RuneCountInString(ent.Message)
		justLen := logger.TermMsgJust - paddingLen
		if len(fields) > 0 {
			if length < justLen {
				line.AppendString(string(bytes.Repeat([]byte{' '}, justLen-length)))
			} else {
				line.AppendString(" -->")
			}
		}
	}

	// write fields
	c.writeContext(line, c.jsonEncoder.files, levelColor)
	c.writeContext(line, fields, levelColor)

	// If there's no stacktrace key, honor that; this allows users to force
	// single-line output.
	if ent.Stack != "" && c.StacktraceKey != "" {
		line.AppendByte('\n')
		line.AppendString(ent.Stack)
	}

	if c.LineEnding != "" {
		line.AppendString(c.LineEnding)
	} else {
		line.AppendString(zapcore.DefaultLineEnding)
	}
	return line, nil
}

func (c *consoleEncoder) isPrint(module string) bool {
	if c.modules == nil || len(c.modules) == 0 {
		return true
	}
	var modulesRegMap sync.Map
	if c.modulesReg.Load() != nil {
		modulesRegMap = c.modulesReg.Load().(sync.Map)
	} else {
		for _, s := range c.modules {
			modulesRegMap.Store(s, s)
		}
		c.modulesReg.Store(modulesRegMap)
	}
	if _, ok := modulesRegMap.Load("*"); ok {
		return true
	}
	if _, ok := modulesRegMap.Load(module); ok {
		return true
	}
	return false
}

func (c *consoleEncoder) writeContext(line *buffer.Buffer, extra []zapcore.Field, levelColor int) {
	for _, field := range extra {
		line.AppendString(" ")
		if levelColor > 0 {
			line.AppendString(logger.ColorRender(field.Key, levelColor, 0) + "=")
		} else {
			line.AppendString(field.Key + "=")
		}
		if len(field.String) > 0 {
			line.AppendString(logger.ToString(field.String))
		} else if field.Interface != nil {
			line.AppendString(logger.ToString(field.Interface))
		} else {
			line.AppendString(logger.ToString(field.Integer))
		}
	}
}

func (c *consoleEncoder) addTabIfNecessary(line *buffer.Buffer) {
	if line.Len() > 0 {
		line.AppendByte('\t')
	}
}

func getColorByLevel(level zapcore.Level) int {
	switch level {
	case zapcore.FatalLevel, zapcore.PanicLevel, zapcore.DPanicLevel:
		return 35
	case zapcore.ErrorLevel:
		return 31
	case zapcore.WarnLevel:
		return 33
	case zapcore.InfoLevel:
		return 32
	case zapcore.DebugLevel:
		return 36
	default:
		return 0
	}
}
