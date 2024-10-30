package customLogger

import (
	"context"
	"github.com/davecgh/go-spew/spew"
	"io"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/apsdehal/go-logger"
)

type LoggerLevel string

const (
	ErrorLevel LoggerLevel = "ERROR"
	WarnLevel  LoggerLevel = "WARN"
	InfoLevel  LoggerLevel = "INFO"
	DebugLevel LoggerLevel = "DEBUG"
)

const (
	typeLog          string = "APACHE"
	envLogLevel      string = "LOG_LEVEL"
	defaultMsgLength int    = 3000
	correlationId    string = "correlationId"
)

type customLog struct {
	opts opts
}

type opts struct {
	//Name of module that used Logger
	module string
	//Config Level. It can be specified using this variable, the global variable LOG_LEVEL or the FilePath configuration file.
	logLevel LoggerLevel
	//Maximum message length. Number of characters
	logMaxLength int
	//Log output
	stdOut io.Writer
	stdErr io.Writer
}

type OptsFunc func(*opts)

func defaultOpts() opts {
	return opts{
		module:       "Logger",
		logLevel:     InfoLevel,
		logMaxLength: defaultMsgLength,
		stdOut:       os.Stdout,
		stdErr:       os.Stderr,
	}
}

func LogLevel(logLevel LoggerLevel) OptsFunc {
	return func(opts *opts) {
		opts.logLevel = logLevel
	}
}

func LogModuleName(name string) OptsFunc {
	return func(opts *opts) {
		opts.module = name
	}
}

func LogMaxLength(numberOfChar int) OptsFunc {
	return func(opts *opts) {
		opts.logMaxLength = numberOfChar
		if numberOfChar <= 0 {
			opts.logMaxLength = defaultMsgLength
		}
	}
}

func StringToLogLevel(level string) LoggerLevel {
	level = strings.ToUpper(level) //Upper string

	switch LoggerLevel(level) {
	case ErrorLevel, WarnLevel, InfoLevel, DebugLevel:
		return LoggerLevel(level)
	default:
		return InfoLevel //default level
	}
}

// NewLog - Logger instantiation
func NewLog(opts ...OptsFunc) Log {
	o := defaultOpts()
	for _, fn := range opts {
		fn(&o)
	}

	if os.Getenv(envLogLevel) != "" { //Get log level from environment variable LOG_LEVEL
		o.logLevel = StringToLogLevel(os.Getenv(envLogLevel))
	}

	return &customLog{opts: o}
}

func (l *customLog) GetTimeNow() time.Time {
	if l.opts.logLevel == DebugLevel {
		return time.Now()
	} else {
		return time.Time{}
	}
}

func (l *customLog) Debug(msg ...any) {
	if l.opts.logLevel == DebugLevel {
		selectLogType(typeLog, l.opts.stdOut, l.opts.module, DebugLevel, l.opts.logMaxLength, nil, msg)
	}
}

func (l *customLog) DebugTime(msg ...any) {
	if l.opts.logLevel == DebugLevel {
		selectLogType(typeLog, l.opts.stdOut, l.opts.module, DebugLevel, l.opts.logMaxLength, nil, msg)
	}
}

func (l *customLog) Info(msg ...any) {
	if l.opts.logLevel == DebugLevel || l.opts.logLevel == InfoLevel {
		selectLogType(typeLog, l.opts.stdOut, l.opts.module, InfoLevel, l.opts.logMaxLength, nil, msg)
	}
}

func (l *customLog) Warning(msg ...any) {
	if l.opts.logLevel == DebugLevel || l.opts.logLevel == InfoLevel || l.opts.logLevel == WarnLevel {
		selectLogType(typeLog, l.opts.stdOut, l.opts.module, WarnLevel, l.opts.logMaxLength, nil, msg)
	}
}

func (l *customLog) Error(msg ...any) {
	selectLogType(typeLog, l.opts.stdErr, l.opts.module, ErrorLevel, l.opts.logMaxLength, nil, msg)
}

// Context
func (l *customLog) DebugCtx(ctx context.Context, msg ...any) {
	if l.opts.logLevel == DebugLevel {
		selectLogType(typeLog, l.opts.stdOut, l.opts.module, DebugLevel, l.opts.logMaxLength, ctx, msg)
	}
}

func (l *customLog) DebugTimeCtx(ctx context.Context, msg ...any) {
	if l.opts.logLevel == DebugLevel {
		selectLogType(typeLog, l.opts.stdOut, l.opts.module, DebugLevel, l.opts.logMaxLength, ctx, msg)
	}
}

func (l *customLog) InfoCtx(ctx context.Context, msg ...any) {
	if l.opts.logLevel == DebugLevel || l.opts.logLevel == InfoLevel {
		selectLogType(typeLog, l.opts.stdOut, l.opts.module, InfoLevel, l.opts.logMaxLength, ctx, msg)
	}
}

func (l *customLog) WarningCtx(ctx context.Context, msg ...any) {
	if l.opts.logLevel == DebugLevel || l.opts.logLevel == InfoLevel || l.opts.logLevel == WarnLevel {
		selectLogType(typeLog, l.opts.stdOut, l.opts.module, WarnLevel, l.opts.logMaxLength, ctx, msg)
	}
}

func (l *customLog) ErrorCtx(ctx context.Context, msg ...any) {
	selectLogType(typeLog, l.opts.stdErr, l.opts.module, ErrorLevel, l.opts.logMaxLength, ctx, msg)
}

func (l *customLog) CalculateDifference(initial time.Time) time.Duration {
	if l.opts.logLevel == DebugLevel {
		return time.Since(initial)
	} else {
		return time.Duration(0)
	}
}

// return context
func (l *customLog) DebugReturnCtx(ctx context.Context, msg ...any) context.Context {
	if l.opts.logLevel == DebugLevel {
		ctx = selectLogType(typeLog, l.opts.stdOut, l.opts.module, DebugLevel, l.opts.logMaxLength, ctx, msg)
	}
	return ctx
}

func (l *customLog) InfoReturnCtx(ctx context.Context, msg ...any) context.Context {
	if l.opts.logLevel == DebugLevel || l.opts.logLevel == InfoLevel {
		ctx = selectLogType(typeLog, l.opts.stdOut, l.opts.module, InfoLevel, l.opts.logMaxLength, ctx, msg)
	}
	return ctx
}

func (l *customLog) WarningReturnCtx(ctx context.Context, msg ...any) context.Context {
	if l.opts.logLevel == DebugLevel || l.opts.logLevel == InfoLevel || l.opts.logLevel == WarnLevel {
		ctx = selectLogType(typeLog, l.opts.stdOut, l.opts.module, WarnLevel, l.opts.logMaxLength, ctx, msg)
	}

	return ctx
}

// Aux functions
func selectLogType(logType string, out io.Writer, module string, level LoggerLevel, logMaxLength int, ctx context.Context, msg []any) context.Context {

	if logType == "APACHE" {
		ctx = apacheLog(out, module, level, generateMsg(msg), logMaxLength, ctx)
	}
	return ctx
}

func apacheLog(out io.Writer, module string, level LoggerLevel, msg string, logMaxLength int, ctx context.Context) context.Context {
	logCustom, err := logger.New(module, out)
	if err != nil {
		panic(err)
	}
	logCustom.SetFormat("[%{time:2006-02-01 15:04:05.000}] [%{module}%{message}")
	logCustom.SetLogLevel(logger.DebugLevel)

	pid := syscall.Getpid()
	//tid := syscall.Gettid()//Ubuntu
	tid := 0
	_, fileName, line, ok := runtime.Caller(3)
	slash := strings.LastIndex(fileName, "/")
	fileName = fileName[slash+1:]
	if !ok {
		fileName = "???"
		line = 0
	}
	var corrId = ""

	if ctx != nil {
		if ctx.Value(correlationId) != nil {
			corrId = ctx.Value(correlationId).(string)
		}
	}

	l := math.Ceil(float64(len(msg)) / float64(logMaxLength))
	length := int(l)
	for i := 0; i < length; i++ {
		var mat string
		if i == length-1 {
			lastBlock := i * logMaxLength
			mat = msg[lastBlock:]
		} else {
			first := i * logMaxLength
			last := first + logMaxLength
			mat = msg[first:last]
		}
		entries := strconv.Itoa(i+1) + "/" + strconv.Itoa(length)

		switch level {
		case ErrorLevel:
			logCustom.ErrorF(":%s][pid:%d tid:%d][correlationId: %s][entries: %s][%s:%d] %s", level, pid, tid, corrId, entries, fileName, line, mat)
		case WarnLevel:
			logCustom.WarningF(":%s][pid:%d tid:%d][correlationId: %s][entries: %s][%s:%d] %s", level, pid, tid, corrId, entries, fileName, line, mat)
		case InfoLevel:
			logCustom.InfoF(":%s][pid:%d tid:%d][correlationId: %s][entries: %s][%s:%d] %s", level, pid, tid, corrId, entries, fileName, line, mat)
		case DebugLevel:
			logCustom.DebugF(":%s][pid:%d tid:%d][correlationId: %s][entries: %s][%s:%d] %s", level, pid, tid, corrId, entries, fileName, line, mat)
		default:
			logCustom.InfoF(":%s][pid:%d tid:%d][correlationId: %s][entries: %s][%s:%d] %s", level, pid, tid, corrId, entries, fileName, line, mat)
		}
	}

	return ctx
}

func generateMsg(msg []any) string {

	msgString := ""

	allNil := true
	for _, arg := range msg {
		if arg != nil {
			allNil = false
			break
		}
	}

	if allNil {
		return msgString
	}

	for i, arg := range msg {
		argString := spew.Sprint(arg)
		msgString += argString
		if i < len(msg)-1 {
			msgString += " "
		}
	}

	return msgString
}
