// Package logrus
//
// @author: xwc1125
// @date: 2021/8/30
package logrus

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"unicode/utf8"

	"github.com/chain5j/logger"
	"github.com/sirupsen/logrus"
)

// Formatter - logrus formatter, implements logrus.Formatter
type Formatter struct {
	// TimestampFormat - default: time.StampMilli = "Jan _2 15:04:05.000"
	TimestampFormat string

	// UseColor - default: disable colors
	UseColor bool

	// TrimMessages - trim whitespaces on messages
	TrimMessages bool

	// LocationEnabled - print caller
	LocationEnabled bool
	LocationTrims   []string
	locationLength  uint32

	// CustomCallerFormatter - set custom formatter for caller info
	CustomCallerFormatter func(*runtime.Frame) string

	Modules    []string
	modulesReg atomic.Value
}

// Format an log entry
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	module := logger.ToString(entry.Data[logger.FieldKeyModule])
	if !f.isPrint(module) {
		return []byte(""), nil
	}
	levelColor := 0
	if f.UseColor {
		levelColor = getColorByLevel(entry.Level)
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = logger.DefaultTermTimestampFormat
	}

	// output buffer
	b := &bytes.Buffer{}

	// write level
	var level = strings.ToUpper(entry.Level.String())
	{
		levelLen := len(level)
		if levelColor > 0 {
			level = logger.ColorRender(level, levelColor, 0)
		}
		b.WriteString(level)
		b.WriteString(strings.Repeat(" ", 5-levelLen))
	}

	// write time
	{
		b.WriteString(" ")
		b.WriteString("[")
		b.WriteString(entry.Time.Format(timestampFormat))
		b.WriteString("]")
	}

	// write module
	{
		if levelColor > 0 {
			module = logger.ColorRender(module, levelColor, 0)
		}
		b.WriteString(" ")
		b.WriteString("[")
		b.WriteString(module)
		b.WriteString("]")
	}

	// write location
	{
		if f.LocationEnabled {
			f.writeCaller(b, entry)
		}
	}

	// write message
	if f.TrimMessages {
		b.WriteString(" ")
		b.WriteString(strings.TrimSpace(entry.Message))
	} else {
		b.WriteString(" ")
		b.WriteString(entry.Message)
	}

	// write padding
	{
		// try to justify the log output for short messages
		length := utf8.RuneCountInString(entry.Message)
		justLen := logger.TermMsgJust - utf8.RuneCountInString(module)
		if len(entry.Data) > 0 {
			if length < justLen {
				b.Write(bytes.Repeat([]byte{' '}, justLen-length))
			} else {
				b.WriteString(" -->")
			}
		}
	}

	// write fields
	{
		f.writeFields(b, entry, levelColor)
	}
	b.WriteByte('\n')

	return b.Bytes(), nil
}

func (f *Formatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		if f.CustomCallerFormatter != nil {
			fmt.Fprintf(b, f.CustomCallerFormatter(entry.Caller))
		} else {
			fmt.Fprintf(
				b,
				" %s:%d",
				entry.Caller.Function,
				entry.Caller.Line,
			)
		}
	}
}

func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry, levelColor int) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			b.WriteString(" ")
			f.writeField(b, entry, field, levelColor)
		}
	}
}

func (f *Formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string, levelColor int) {
	if !strings.EqualFold(field, logger.FieldKeyModule) {
		if levelColor > 0 {
			fmt.Fprintf(b, "%s", logger.ColorRender(field, levelColor, 0))
		} else {
			fmt.Fprintf(b, "%s", field)
		}
		fmt.Fprintf(b, "=%v", entry.Data[field])
	}
}

func getColorByLevel(level logrus.Level) int {
	switch level {
	case logrus.FatalLevel, logrus.PanicLevel:
		return 35
	case logrus.ErrorLevel:
		return 31
	case logrus.WarnLevel:
		return 33
	case logrus.InfoLevel:
		return 32
	case logrus.DebugLevel:
		return 36
	case logrus.TraceLevel:
		return 34
	default:
		return 0
	}
}

func (f *Formatter) isPrint(module string) bool {
	if f.Modules == nil || len(f.Modules) == 0 {
		return true
	}
	var modulesRegMap sync.Map
	if f.modulesReg.Load() != nil {
		modulesRegMap = f.modulesReg.Load().(sync.Map)
	} else {
		for _, s := range f.Modules {
			modulesRegMap.Store(s, s)
		}
		f.modulesReg.Store(modulesRegMap)
	}
	if _, ok := modulesRegMap.Load("*"); ok {
		return true
	}
	if _, ok := modulesRegMap.Load(module); ok {
		return true
	}
	return false
}
