package logutil

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/justmumu/goutils/fileutil"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Configuration for logging
type LoggerConfig struct {
	// Name is the name of the logger
	Name string

	ConsoleEnabled bool
	ConsoleLevel   LogLevel
	ConsoleJson    bool

	FileEnabled bool
	FileLevel   LogLevel
	FileJson    bool

	// LogDirectory to log to when file logging is enabled
	LogDirectory string
	// Filename is the name of the log file which will be placed inside the directory
	Filename string
	// MaxSize the max size in MB of the logfile before it is rotated
	MaxSize int
	// MaxBackups the max number of rotated files to keep
	MaxBackup int
	// MaxAge the max age in days to keep a logfile
	MaxAge int
}

type Logger interface {
	// Debug logs the provided arguments at [DebugLevel]. Spaces are added between arguments when neither is a string.
	Debug(args ...interface{})
	// Debugf formats the message according to the format specifier and logs it at [DebugLevel].
	Debugf(template string, args ...interface{})
	// Debugln logs a message at [DebugLevel]. Spaces are always added between arguments.
	Debugln(args ...interface{})
	// Debugw logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
	//
	// When debug-level logging is disabled, this is much faster than
	//	s.With(keysAndValues).Debug(msg)
	Debugw(msg string, keysAndValues ...interface{})

	// Info logs the provided arguments at [InfoLevel]. Spaces are added between arguments when neither is a string.
	Info(args ...interface{})
	// Infof formats the message according to the format specifier and logs it at [InfoLevel].
	Infof(template string, args ...interface{})
	// Infoln logs a message at [InfoLevel]. Spaces are always added between arguments.
	Infoln(args ...interface{})
	// Infow logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
	Infow(msg string, keysAndValues ...interface{})

	// Warn logs the provided arguments at [WarnLevel]. Spaces are added between arguments when neither is a string.
	Warn(args ...interface{})
	// Warnf formats the message according to the format specifier
	// and logs it at [WarnLevel].
	Warnf(template string, args ...interface{})
	// Warnln logs a message at [WarnLevel].
	// Spaces are always added between arguments.
	Warnln(args ...interface{})
	// Warnw logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	Warnw(msg string, keysAndValues ...interface{})

	// Error logs the provided arguments at [ErrorLevel].
	// Spaces are added between arguments when neither is a string.
	Error(args ...interface{})
	// Errorf formats the message according to the format specifier
	// and logs it at [ErrorLevel].
	Errorf(template string, args ...interface{})
	// Errorln logs a message at [ErrorLevel].
	// Spaces are always added between arguments.
	Errorln(args ...interface{})
	// Errorw logs a message with some additional context. The variadic key-value
	// pairs are treated as they are in With.
	Errorw(msg string, keysAndValues ...interface{})

	// Panic constructs a message with the provided arguments and panics.
	// Spaces are added between arguments when neither is a string.
	Panic(args ...interface{})
	// Panicf formats the message according to the format specifier
	// and panics.
	Panicf(template string, args ...interface{})
	// Panicln logs a message at [PanicLevel] and panics.
	// Spaces are always added between arguments.
	Panicln(args ...interface{})
	// Panicw logs a message with some additional context, then panics. The
	// variadic key-value pairs are treated as they are in With.
	Panicw(msg string, keysAndValues ...interface{})

	// Fatal constructs a message with the provided arguments and calls os.Exit.
	// Spaces are added between arguments when neither is a string.
	Fatal(args ...interface{})
	// Fatalf formats the message according to the format specifier
	// and calls os.Exit.
	Fatalf(template string, args ...interface{})
	// Fatalln logs a message at [FatalLevel] and calls os.Exit.
	// Spaces are always added between arguments.
	Fatalln(args ...interface{})
	// Fatalw logs a message with some additional context, then calls os.Exit. The
	// variadic key-value pairs are treated as they are in With.
	Fatalw(msg string, keysAndValues ...interface{})

	// Named adds a new path segment to the logger's name. Segments are joined by
	// periods. By default, Loggers are unnamed.
	Named(s string) Logger
	// Sync calls the underlying Core's Sync method, flushing any buffered log
	// entries. Applications should take care to call Sync before exiting.
	Sync() error
	// SetConsoleLevel sets the logging level for the console logger.
	SetConsoleLevel(level LogLevel)
	// SetFileLevel sets the logging level for the file logger.
	SetFileLevel(level LogLevel)
}

type logger struct {
	consoleAtomLvl zap.AtomicLevel
	fileAtomLvl    zap.AtomicLevel

	unsugared *zap.Logger
	*zap.SugaredLogger
}

func NewLogger(config LoggerConfig) Logger {
	ll := &logger{}
	// Prepare logging level
	var consoleLevel zapcore.Level
	consoleLevel.Set(strings.ToLower(config.ConsoleLevel.String()))
	ll.consoleAtomLvl = zap.NewAtomicLevelAt(consoleLevel)

	var fileLevel zapcore.Level
	fileLevel.Set(strings.ToLower(config.FileLevel.String()))
	ll.fileAtomLvl = zap.NewAtomicLevelAt(fileLevel)

	// Prepare encoder configs
	consoleEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
	}

	jsonEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.NanosDurationEncoder,
	}

	// Prepare encoders
	jsonEncoder := zapcore.NewJSONEncoder(jsonEncoderConfig)

	// Prepare zap cores
	var cores []zapcore.Core

	if config.ConsoleEnabled {
		if config.ConsoleJson {
			cores = append(cores, zapcore.NewCore(jsonEncoder, zapcore.Lock(os.Stderr), ll.consoleAtomLvl))
		} else {
			coloredTextEncoderConfig := consoleEncoderConfig
			coloredTextEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			coloredTextEncoder := zapcore.NewConsoleEncoder(coloredTextEncoderConfig)
			cores = append(cores, zapcore.NewCore(coloredTextEncoder, zapcore.Lock(os.Stderr), ll.consoleAtomLvl))
		}
	}

	if config.FileEnabled {
		fileSyncer := newRotateFile(config)
		if fileSyncer != nil {
			if config.FileJson {
				cores = append(cores, zapcore.NewCore(jsonEncoder, newRotateFile(config), ll.fileAtomLvl))
			} else {
				uncoloredTextEncoderConfig := consoleEncoderConfig
				uncoloredTextEncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
				uncoloredTextEncoder := zapcore.NewConsoleEncoder(uncoloredTextEncoderConfig)
				cores = append(cores, zapcore.NewCore(uncoloredTextEncoder, newRotateFile(config), ll.fileAtomLvl))
			}
		}
	}
	core := zapcore.NewTee(cores...)

	// Prepare zap logger instance
	unsugared := zap.New(core)

	if config.Name != "" {
		unsugared = unsugared.Named(strings.ReplaceAll(strings.ToLower(config.Name), " ", "_"))
	}
	ll.unsugared = unsugared
	ll.SugaredLogger = unsugared.Sugar()

	return ll
}

// SetConsoleLevel sets the logging level for the console logger.
func (l *logger) SetConsoleLevel(level LogLevel) {
	var consoleLevel zapcore.Level
	consoleLevel.Set(strings.ToLower(level.String()))
	l.consoleAtomLvl.SetLevel(consoleLevel)
}

// SetFileLevel sets the logging level for the file logger.
func (l *logger) SetFileLevel(level LogLevel) {
	var fileLevel zapcore.Level
	fileLevel.Set(strings.ToLower(level.String()))
	l.fileAtomLvl.SetLevel(fileLevel)
}

// Named adds a new path segment to the logger's name. Segments are joined by periods. By default, Loggers are unnamed.
func (l *logger) Named(s string) Logger {
	ll := l.unsugared.Named(s)
	return &logger{
		consoleAtomLvl: l.consoleAtomLvl,
		fileAtomLvl:    l.fileAtomLvl,
		unsugared:      ll,
		SugaredLogger:  ll.Sugar(),
	}
}

// Sync calls the underlying loggers's Sync method, flushing any buffered log entries. Applications should take care to call Sync before exiting.
func (l *logger) Sync() error {
	return l.unsugared.Sync()
}

func newRotateFile(config LoggerConfig) zapcore.WriteSyncer {
	if err := fileutil.CreateFolders(config.LogDirectory); err != nil {
		Errorf("could not create log directory. err: %v", err)
		return nil
	}

	// Lumberjack.Logger is already safe for concurrent use, so we don't need to lock it.
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(config.LogDirectory, config.Filename),
		MaxSize:    config.MaxSize,
		MaxAge:     config.MaxAge,
		MaxBackups: config.MaxBackup,
	})
}
