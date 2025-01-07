package log

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-stack/stack"
)

const timeKey = "t"
const lvlKey = "lvl"
const msgKey = "msg"
const ctxKey = "ctx"
const errorKey = "LOG15_ERROR"
const skipLevel = 2

type Lvl int

const (
	LvlCrit Lvl = iota
	LvlError
	LvlWarn
	LvlInfo
	LvlDebug
	LvlTrace
)

const ModuleKey = "module"

// AlignedString returns a 5-character string containing the name of a Lvl.
func (l Lvl) AlignedString() string {
	switch l {
	case LvlTrace:
		return "TRACE"
	case LvlDebug:
		return "DEBUG"
	case LvlInfo:
		return "INFO "
	case LvlWarn:
		return "WARN "
	case LvlError:
		return "ERROR"
	case LvlCrit:
		return "CRIT "
	default:
		panic("bad level")
	}
}

// Strings returns the name of a Lvl.
func (l Lvl) String() string {
	switch l {
	case LvlTrace:
		return "trce"
	case LvlDebug:
		return "dbug"
	case LvlInfo:
		return "info"
	case LvlWarn:
		return "warn"
	case LvlError:
		return "eror"
	case LvlCrit:
		return "crit"
	default:
		panic("bad level")
	}
}

// LvlFromString returns the appropriate Lvl from a string name.
// Useful for parsing command line args and configuration files.
func LvlFromString(lvlString string) (Lvl, error) {
	switch lvlString {
	case "trace", "trce":
		return LvlTrace, nil
	case "debug", "dbug":
		return LvlDebug, nil
	case "info":
		return LvlInfo, nil
	case "warn":
		return LvlWarn, nil
	case "error", "eror":
		return LvlError, nil
	case "crit":
		return LvlCrit, nil
	default:
		return LvlDebug, fmt.Errorf("Unknown level: %v", lvlString)
	}
}

// A Record is what a Logger asks its handler to write
type Record struct {
	Module   string
	Time     time.Time
	Lvl      Lvl
	Msg      string
	Ctx      []interface{}
	Call     stack.Call
	KeyNames RecordKeyNames
}

// RecordKeyNames gets stored in a Record when the write function is executed.
type RecordKeyNames struct {
	Time   string
	Lvl    string
	Module string
	Msg    string
	Ctx    string
}

// A Logger writes key/value pairs to a Handler
type Logger interface {
	// New returns a new Logger that has this logger's context plus the given context
	New(module string, ctx ...interface{}) Logger

	// GetHandler gets the handler associated with the logger.
	GetHandler() Handler

	// SetHandler updates the logger to write records to the specified handler.
	SetHandler(h Handler)

	// Log a message at the given level with context key/value pairs
	Trace(msg string, ctx ...interface{})
	Debug(msg string, ctx ...interface{})
	Info(msg string, ctx ...interface{})
	Warn(msg string, ctx ...interface{})
	Error(msg string, ctx ...interface{})
	Crit(msg string, ctx ...interface{})

	// log
	Printf(format string, v ...interface{})
	Print(v ...interface{})
	Println(v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
}

type logger struct {
	*log.Logger
	module string
	ctx    []interface{}
	h      *swapHandler
}

func (l *logger) write(msg string, lvl Lvl, ctx []interface{}, skip int) {
	l.h.Log(&Record{
		Module: l.module,
		Time:   time.Now(),
		Lvl:    lvl,
		Msg:    msg,
		Ctx:    newContext(l.ctx, ctx),
		Call:   stack.Caller(skip),
		KeyNames: RecordKeyNames{
			Time:   timeKey,
			Lvl:    lvlKey,
			Module: ModuleKey,
			Msg:    msgKey,
			Ctx:    ctxKey,
		},
	})
}

var std = log.New(os.Stderr, "", log.LstdFlags)

func (l *logger) New(module string, ctx ...interface{}) Logger {
	child := &logger{
		Logger: std,
		module: module,
		ctx:    newContext(l.ctx, ctx),
		h:      new(swapHandler),
	}
	child.SetHandler(l.h)
	return child
}

func newContext(prefix []interface{}, suffix []interface{}) []interface{} {
	normalizedSuffix := normalize(suffix)
	newCtx := make([]interface{}, len(prefix)+len(normalizedSuffix))
	n := copy(newCtx, prefix)
	copy(newCtx[n:], normalizedSuffix)
	return newCtx
}

func (l *logger) Trace(msg string, ctx ...interface{}) {
	l.write(msg, LvlTrace, ctx, skipLevel)
}

func (l *logger) Debug(msg string, ctx ...interface{}) {
	l.write(msg, LvlDebug, ctx, skipLevel)
}

func (l *logger) Info(msg string, ctx ...interface{}) {
	l.write(msg, LvlInfo, ctx, skipLevel)
}

func (l *logger) Warn(msg string, ctx ...interface{}) {
	l.write(msg, LvlWarn, ctx, skipLevel)
}

func (l *logger) Error(msg string, ctx ...interface{}) {
	l.write(msg, LvlError, ctx, skipLevel)
}

func (l *logger) Crit(msg string, ctx ...interface{}) {
	l.write(msg, LvlCrit, ctx, skipLevel)
	os.Exit(1)
}

// native

func (l *logger) Printf(format string, v ...interface{}) {
	l.write(fmt.Sprintf(format, v...), LvlInfo, nil, skipLevel)
}

func (l *logger) Print(v ...interface{}) {
	l.write(fmt.Sprint(v...), LvlInfo, nil, skipLevel)
}

func (l *logger) Println(v ...interface{}) {
	l.write(fmt.Sprintln(v...), LvlInfo, nil, skipLevel)
}

func (l *logger) Fatal(v ...interface{}) {
	l.write(fmt.Sprint(v...), LvlCrit, nil, skipLevel)
	os.Exit(1)
}

func (l *logger) Fatalf(format string, v ...interface{}) {
	l.write(fmt.Sprintf(format, v...), LvlCrit, nil, skipLevel)
	os.Exit(1)
}

func (l *logger) Fatalln(v ...interface{}) {
	l.write(fmt.Sprintln(v...), LvlCrit, nil, skipLevel)
	os.Exit(1)
}

func (l *logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.write(s, LvlCrit, nil, skipLevel)
	panic(s)
}

func (l *logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.write(s, LvlCrit, nil, skipLevel)
	panic(s)
}

func (l *logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.write(s, LvlCrit, nil, skipLevel)
	panic(s)
}

func (l *logger) GetHandler() Handler {
	return l.h.Get()
}

func (l *logger) SetHandler(h Handler) {
	l.h.Swap(h)
}

func normalize(ctx []interface{}) []interface{} {
	// if the caller passed a Ctx object, then expand it
	if len(ctx) == 1 {
		if ctxMap, ok := ctx[0].(Ctx); ok {
			ctx = ctxMap.toArray()
		}
	}

	// ctx needs to be even because it's a series of key/value pairs
	// no one wants to check for errors on logging functions,
	// so instead of erroring on bad input, we'll just make sure
	// that things are the right length and users can fix bugs
	// when they see the output looks wrong
	if len(ctx)%2 != 0 {
		ctx = append(ctx, nil, errorKey, "Normalized odd number of arguments by adding nil")
	}

	return ctx
}

// Lazy allows you to defer calculation of a logged value that is expensive
// to compute until it is certain that it must be evaluated with the given filters.
//
// Lazy may also be used in conjunction with a Logger's New() function
// to generate a child logger which always reports the current value of changing
// state.
//
// You may wrap any function which takes no arguments to Lazy. It may return any
// number of values of any type.
type Lazy struct {
	Fn interface{}
}

// Ctx is a map of key/value pairs to pass as context to a log function
// Use this only if you really need greater safety around the arguments you pass
// to the logging functions.
type Ctx map[string]interface{}

func (c Ctx) toArray() []interface{} {
	arr := make([]interface{}, len(c)*2)

	i := 0
	for k, v := range c {
		arr[i] = k
		arr[i+1] = v
		i += 2
	}

	return arr
}
