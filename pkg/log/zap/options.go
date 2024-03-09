package zap

import (
	"github.com/youngfs/youngfs/pkg/log"
	"go.uber.org/zap/zapcore"
	"io"
)

type config struct {
	service  string
	path     string
	level    log.Level
	zapLevel zapcore.Level
	debug    bool
	maxSize  int
	maxAge   int
	writer   io.Writer
}

func (cfg *config) updateLevel(lvl log.Level) {
	cfg.level = lvl
	cfg.zapLevel = levelToZapLevel[lvl]
}

func defaultCfg(service string) *config {
	return &config{
		service:  service,
		level:    log.InfoLevel,
		zapLevel: zapcore.InfoLevel,
		debug:    false,
		maxSize:  32,
		maxAge:   0,
	}
}

type Option interface {
	apply(*config)
}

type optionFunc func(*config)

func (f optionFunc) apply(cfg *config) {
	f(cfg)
}

// WithLevel returns an Option to set the logging level for the Logger.
// The provided lvl will determine the verbosity of logging output.
func WithLevel(lvl log.Level) Option {
	return optionFunc(func(cfg *config) {
		cfg.updateLevel(lvl)
	})
}

// WithLogFileSize returns an Option to set the maximum file size for the Logger's log files.
// The provided size is in MiB (Mebibytes) and determines the maximum size a log file can grow
// before it gets rotated.
func WithLogFileSize(size int) Option {
	return optionFunc(func(cfg *config) {
		cfg.maxSize = size
	})
}

// WithLogFileAge returns an Option to set the maximum age for the Logger's log files.
// The provided age is in days and determines the maximum number of days a log file is retained.
// If age is set to 0, log files will not be deleted based on age.
func WithLogFileAge(age int) Option {
	return optionFunc(func(cfg *config) {
		cfg.maxAge = age
	})
}

// WithLogPath returns an Option to set the directory path where the Logger's log files will be saved.
// The provided path determines the destination folder for all log files created by the Logger.
func WithLogPath(path string) Option {
	return optionFunc(func(cfg *config) {
		cfg.path = path
	})
}

// WithDebug returns an Option to configure the Logger for a debug environment.
// When this option is set:
// 1. Log files will not be saved locally.
// 2. Log messages will be output to the terminal.
// 3. The log level will be set to Debug.
// This is useful to ensure that unit tests do not produce persistent log files
// and provide verbose logging output for debugging.
func WithDebug() Option {
	return optionFunc(func(cfg *config) {
		cfg.debug = true
		cfg.updateLevel(log.DebugLevel)
	})
}

func WithLogWriter(writer io.Writer) Option {
	return optionFunc(func(cfg *config) {
		cfg.writer = writer
	})
}
