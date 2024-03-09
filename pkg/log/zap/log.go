package zap

import (
	"github.com/natefinch/lumberjack"
	"github.com/youngfs/youngfs/pkg/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
)

var levelToZapLevel = map[log.Level]zapcore.Level{
	log.DebugLevel:  zap.DebugLevel,
	log.InfoLevel:   zap.InfoLevel,
	log.WarnLevel:   zap.WarnLevel,
	log.ErrorLevel:  zap.ErrorLevel,
	log.DPanicLevel: zap.DPanicLevel,
	log.PanicLevel:  zap.PanicLevel,
}

const (
	timeFormat    = "2006-01-02 15:04:05.000000" //time.RFC3339Nano
	logExtension  = ".log"
	timeKey       = "time"
	levelKey      = "level"
	nameKey       = "logger"
	callerKey     = "caller"
	messageKey    = "message"
	stacktraceKey = "stacktrace"
	serviceKey    = "service"
)

var encoderConfig = zapcore.EncoderConfig{
	TimeKey:        timeKey,
	LevelKey:       levelKey,
	NameKey:        nameKey,
	CallerKey:      callerKey,
	MessageKey:     messageKey,
	StacktraceKey:  stacktraceKey,
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.TimeEncoderOfLayout(timeFormat),
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
	EncodeName:     zapcore.FullNameEncoder,
}

func (cfg *config) newWriter(name string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:  path.Join(cfg.path, name),
		MaxSize:   cfg.maxSize,
		MaxAge:    cfg.maxAge,
		LocalTime: true,
		Compress:  false,
	}
	var writers []zapcore.WriteSyncer
	if cfg.writer == nil {
		if !cfg.debug {
			writers = append(writers, zapcore.AddSync(lumberJackLogger))
		}
		if cfg.level == log.DebugLevel {
			writers = append(writers, zapcore.AddSync(os.Stdout))
		}
	} else {
		writers = append(writers, zapcore.AddSync(cfg.writer))
	}
	multiWriter := zapcore.NewMultiWriteSyncer(writers...)
	return multiWriter
}

func New(service string, options ...Option) *zap.SugaredLogger {
	cfg := defaultCfg(service)
	for _, opt := range options {
		opt.apply(cfg)
	}

	var cores []zapcore.Core
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// error log
	{
		maxLevel := max(zapcore.ErrorLevel, cfg.zapLevel)
		enableLevel := zap.NewAtomicLevelAt(maxLevel)
		core := zapcore.NewCore(encoder, cfg.newWriter(path.Join("error", cfg.service+logExtension)), enableLevel)
		cores = append(cores, core)
	}

	// cfg log
	if cfg.level < log.ErrorLevel {
		enableLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= cfg.zapLevel && lvl < zap.ErrorLevel
		})
		core := zapcore.NewCore(encoder, cfg.newWriter(path.Join("info", cfg.service+logExtension)), enableLevel)
		cores = append(cores, core)
	}

	core := zapcore.NewTee(cores...)
	var zapOptions []zap.Option
	zapOptions = append(zapOptions, zap.AddCaller()) // stack trace: show file name and line number
	if cfg.level == log.DebugLevel {
		zapOptions = append(zapOptions, zap.Development()) // development mode: makes DPanic-level logs panic instead of simply logging an error
	}
	zapOptions = append(zapOptions, zap.Fields(zap.String(serviceKey, service))) // set filer "service" : ServiceName

	zapLogger := zap.New(core, zapOptions...)
	return zapLogger.Sugar()
}
