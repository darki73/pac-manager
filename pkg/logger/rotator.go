package logger

import (
	"io"

	logrusPackage "github.com/sirupsen/logrus"
	lumberjackPackage "gopkg.in/natefinch/lumberjack.v2"
)

// RotatorFileConfig represents configuration for log rotation of a file
type RotatorFileConfig struct {
	// FileName defines output file name
	FileName string
	// MaxSize defines maximum log size in megabytes
	MaxSize uint
	// MaxBackups defines maximum number of backups for one file
	MaxBackups uint
	// MaxAge defines maximum age of log in days
	MaxAge uint
	// Level defines log level
	Level logrusPackage.Level
	// Formatter defines formatter used by rotator
	Formatter logrusPackage.Formatter
	// Compress defines whether logs should be compressed
	Compress bool
}

// OutputRotatorFileHook represents structure of `output` file hook
type OutputRotatorFileHook struct {
	// Configuration defines rotator configuration
	Configuration RotatorFileConfig
	// writer defines writer used by rotator
	writer io.Writer
}

// NewOutputRotatorFileHook creates new instance of rotator file hook (STDOUT)
func NewOutputRotatorFileHook(config RotatorFileConfig) logrusPackage.Hook {
	hook := OutputRotatorFileHook{
		Configuration: config,
	}
	hook.writer = &lumberjackPackage.Logger{
		Filename:   config.FileName,
		MaxSize:    int(config.MaxSize),
		MaxBackups: int(config.MaxBackups),
		MaxAge:     int(config.MaxAge),
		Compress:   config.Compress,
	}
	return &hook
}

// Levels returns list of levels which should be processed by the rotator hook
func (hook *OutputRotatorFileHook) Levels() []logrusPackage.Level {
	return []logrusPackage.Level{
		logrusPackage.InfoLevel,
		logrusPackage.DebugLevel,
		logrusPackage.WarnLevel,
		logrusPackage.TraceLevel,
	}
}

// Fire fires writer event
func (hook *OutputRotatorFileHook) Fire(entry *logrusPackage.Entry) error {
	bytes, err := hook.Configuration.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

// ErrorRotatorFileHook represents structure of `error` file hook
type ErrorRotatorFileHook struct {
	// Configuration defines rotator configuration
	Configuration RotatorFileConfig
	// writer defines writer used by rotator
	writer io.Writer
}

// NewErrorRotatorFileHook creates new instance of rotator file hook (STDERR)
func NewErrorRotatorFileHook(config RotatorFileConfig) logrusPackage.Hook {
	hook := ErrorRotatorFileHook{
		Configuration: config,
	}
	hook.writer = &lumberjackPackage.Logger{
		Filename:   config.FileName,
		MaxSize:    int(config.MaxSize),
		MaxBackups: int(config.MaxBackups),
		MaxAge:     int(config.MaxAge),
		Compress:   config.Compress,
	}
	return &hook
}

// Levels returns list of levels which should be processed by the rotator hook
func (hook *ErrorRotatorFileHook) Levels() []logrusPackage.Level {
	return []logrusPackage.Level{
		logrusPackage.PanicLevel,
		logrusPackage.FatalLevel,
		logrusPackage.ErrorLevel,
	}
}

// Fire fires writer event
func (hook *ErrorRotatorFileHook) Fire(entry *logrusPackage.Entry) error {
	bytes, err := hook.Configuration.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
