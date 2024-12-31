package logutil

var (
	DefaultLogger Logger
)

func init() {
	DefaultLogger = NewLogger(LoggerConfig{
		ConsoleEnabled: true,
		ConsoleLevel:   InfoLevel,
		ConsoleJson:    false,
		FileEnabled:    false,
		FileJson:       true,
	})
}

// Debug logs the provided arguments at [DebugLevel]. Spaces are added between arguments when neither is a string.
func Debug(args ...interface{}) {
	DefaultLogger.Debug(args...)
}

// Debugf formats the message according to the format specifier and logs it at [DebugLevel].
func Debugf(template string, args ...interface{}) {
	DefaultLogger.Debugf(template, args...)
}

// Debugln logs a message at [DebugLevel]. Spaces are always added between arguments.
func Debugln(args ...interface{}) {
	DefaultLogger.Debugln(args...)
}

// Debugw logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//
//	s.With(keysAndValues).Debug(msg)
func Debugw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Debugw(msg, keysAndValues...)
}

// Info logs the provided arguments at [InfoLevel]. Spaces are added between arguments when neither is a string.
func Info(args ...interface{}) {
	DefaultLogger.Info(args...)
}

// Infof formats the message according to the format specifier and logs it at [InfoLevel].
func Infof(template string, args ...interface{}) {
	DefaultLogger.Infof(template, args...)
}

// Infoln logs a message at [InfoLevel]. Spaces are always added between arguments.
func Infoln(args ...interface{}) {
	DefaultLogger.Infoln(args...)
}

// Infow logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Infow(msg, keysAndValues...)
}

// Warn logs the provided arguments at [WarnLevel]. Spaces are added between arguments when neither is a string.
func Warn(args ...interface{}) {
	DefaultLogger.Warn(args...)
}

// Warnf formats the message according to the format specifier
// and logs it at [WarnLevel].
func Warnf(template string, args ...interface{}) {
	DefaultLogger.Warnf(template, args...)
}

// Warnln logs a message at [WarnLevel].
// Spaces are always added between arguments.
func Warnln(args ...interface{}) {
	DefaultLogger.Warnln(args...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Warnw(msg, keysAndValues...)
}

// Error logs the provided arguments at [ErrorLevel].
// Spaces are added between arguments when neither is a string.
func Error(args ...interface{}) {
	DefaultLogger.Error(args...)
}

// Errorf formats the message according to the format specifier
// and logs it at [ErrorLevel].
func Errorf(template string, args ...interface{}) {
	DefaultLogger.Errorf(template, args...)
}

// Errorln logs a message at [ErrorLevel].
// Spaces are always added between arguments.
func Errorln(args ...interface{}) {
	DefaultLogger.Errorln(args...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Errorw(msg, keysAndValues...)
}

// Panic constructs a message with the provided arguments and panics.
// Spaces are added between arguments when neither is a string.
func Panic(args ...interface{}) {
	DefaultLogger.Panic(args...)
}

// Panicf formats the message according to the format specifier
// and panics.
func Panicf(template string, args ...interface{}) {
	DefaultLogger.Panicf(template, args...)
}

// Panicln logs a message at [PanicLevel] and panics.
// Spaces are always added between arguments.
func Panicln(args ...interface{}) {
	DefaultLogger.Panicln(args...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Panicw(msg, keysAndValues...)
}

// Fatal constructs a message with the provided arguments and calls os.Exit.
// Spaces are added between arguments when neither is a string.
func Fatal(args ...interface{}) {
	DefaultLogger.Fatal(args...)
}

// Fatalf formats the message according to the format specifier
// and calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	DefaultLogger.Fatalf(template, args...)
}

// Fatalln logs a message at [FatalLevel] and calls os.Exit.
// Spaces are always added between arguments.
func Fatalln(args ...interface{}) {
	DefaultLogger.Fatalln(args...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keysAndValues ...interface{}) {
	DefaultLogger.Fatalw(msg, keysAndValues...)
}

// Named adds a new path segment to the logger's name. Segments are joined by
// periods. By default, Loggers are unnamed.
func Named(s string) Logger {
	return DefaultLogger.Named(s)
}

// Sync calls the underlying Core's Sync method, flushing any buffered log
// entries. Applications should take care to call Sync before exiting.
func Sync() error {
	return DefaultLogger.Sync()
}

// SetConsoleLevel sets the console log level
func SetConsoleLevel(level LogLevel) {
	DefaultLogger.SetConsoleLevel(level)
}

// SetFileLevel sets the file log level
func SetFileLevel(level LogLevel) {
	DefaultLogger.SetFileLevel(level)
}
