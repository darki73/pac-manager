package logger

import (
	"runtime"
	"time"

	versionPackage "github.com/darki73/pac-manager/pkg/version"
	colorablePackage "github.com/mattn/go-colorable"
	logrusPackage "github.com/sirupsen/logrus"
)

var (
	// version represents binary version
	version = ""

	// commit represents commit with which binary was built
	commit = ""
)

// init handles package initialization
func init() {
	version = versionPackage.GetVersion()
	commit = versionPackage.GetCommit()
	logrusPackage.SetOutput(colorablePackage.NewColorableStdout())
	logrusPackage.SetFormatter(&logrusPackage.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})
}

// defaultLogMethod default logging method used by the package
func defaultLogMethod(source string) *logrusPackage.Entry {
	return logrusPackage.WithFields(logrusPackage.Fields{
		"version":      version,
		"commit":       commit,
		"source":       source,
		"architecture": runtime.GOARCH,
		"runtime":      runtime.Version(),
	})
}

// SetLogLevelFromString sets log level from string
func SetLogLevelFromString(levelAsString string) {
	level := logrusPackage.ErrorLevel
	switch levelAsString {
	case "t", "trace", "Trace", "TRACE":
		level = logrusPackage.TraceLevel
	case "d", "debug", "Debug", "DEBUG":
		level = logrusPackage.DebugLevel
	case "i", "info", "Info", "INFO":
		level = logrusPackage.InfoLevel
	case "w", "warn", "Warn", "WARN":
		level = logrusPackage.WarnLevel
	case "e", "error", "Error", "ERROR":
		level = logrusPackage.ErrorLevel
	case "f", "fatal", "Fatal", "FATAL":
		level = logrusPackage.FatalLevel
	case "p", "panic", "Panic", "PANIC":
		level = logrusPackage.PanicLevel
	default:
		level = logrusPackage.ErrorLevel
	}
	SetLogLevel(level)
}

// SetLogLevel sets log level from a list of valid logrus levels
func SetLogLevel(level logrusPackage.Level) {
	logrusPackage.SetLevel(level)
}

// GetLogLevel return current log level as a valid logrus level
func GetLogLevel() logrusPackage.Level {
	return logrusPackage.GetLevel()
}

// GetLogLevelAsString return current log level as string
func GetLogLevelAsString() string {
	return GetLogLevel().String()
}

// Trace prints out message in the Trace type logger
func Trace(source, message string) {
	defaultLogMethod(source).Trace(message)
}

// Tracef allows to print out message in the Trace type logger with given format
func Tracef(source, format string, args ...interface{}) {
	defaultLogMethod(source).Tracef(format, args...)
}

// Debug prints out message in the Debug type logger
func Debug(source, message string) {
	defaultLogMethod(source).Debug(message)
}

// Debugf allows to print out message in the Debug type logger with given format
func Debugf(source, format string, args ...interface{}) {
	defaultLogMethod(source).Debugf(format, args...)
}

// Info prints out message in the Info type logger
func Info(source, message string) {
	defaultLogMethod(source).Info(message)
}

// Infof allows to print out message in the Info type logger with given format
func Infof(source, format string, args ...interface{}) {
	defaultLogMethod(source).Infof(format, args...)
}

// Warn prints out message in the Warn type logger
func Warn(source, message string) {
	defaultLogMethod(source).Warn(message)
}

// Warnf allows to print out message in the Warn type logger with given format
func Warnf(source, format string, args ...interface{}) {
	defaultLogMethod(source).Warnf(format, args...)
}

// Error prints out message in the Error type logger
func Error(source, message string) {
	defaultLogMethod(source).Error(message)
}

// Errorf allows to print out message in the Error type logger with given format
func Errorf(source, format string, args ...interface{}) {
	defaultLogMethod(source).Errorf(format, args...)
}

// Fatal prints out message in the Fatal type logger
func Fatal(source, message string) {
	defaultLogMethod(source).Fatal(message)
}

// Fatalf allows to print out message in the Fatal type logger with given format
func Fatalf(source, format string, args ...interface{}) {
	defaultLogMethod(source).Fatalf(format, args...)
}

// Panic prints out message in the Panic type logger
func Panic(source, message string) {
	defaultLogMethod(source).Panic(message)
}

// Panicf allows to print out message in the Panic type logger with given format
func Panicf(source, format string, args ...interface{}) {
	defaultLogMethod(source).Panicf(format, args...)
}

// RegisterRotators register both `output` and `error` log rotators
func RegisterRotators(outputRotatorConfig RotatorFileConfig, errorRotatorConfig RotatorFileConfig) {
	RegisterOutputRotator(outputRotatorConfig)
	RegisterErrorRotator(errorRotatorConfig)
}

// RegisterOutputRotator register `output` log rotator
func RegisterOutputRotator(config RotatorFileConfig) {
	rotatorHook := NewOutputRotatorFileHook(config)
	logrusPackage.AddHook(rotatorHook)
}

// RegisterErrorRotator register `error` log rotator
func RegisterErrorRotator(config RotatorFileConfig) {
	rotatorHook := NewErrorRotatorFileHook(config)
	logrusPackage.AddHook(rotatorHook)
}
