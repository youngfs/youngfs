package log

import (
	"github.com/natefinch/lumberjack"
	"github.com/youngfs/youngfs/vars"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.SugaredLogger

func InitLogger() {

	encoderConfig := zapcore.EncoderConfig{
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

	encoder := zapcore.NewJSONEncoder(encoderConfig) // zapcore.NewConsoleEncoder(encoderConfig)

	lumberJackLogger := &lumberjack.Logger{
		Filename:  vars.ServerName + ".log",
		MaxSize:   16,
		LocalTime: true,
		Compress:  false,
	}

	atomicLevel := zap.NewAtomicLevelAt(zap.InfoLevel) // log level
	if !vars.InfoLog {
		atomicLevel = zap.NewAtomicLevelAt(zap.ErrorLevel)
	}
	if vars.Debug {
		atomicLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	var writers []zapcore.WriteSyncer
	writers = append(writers, zapcore.AddSync(lumberJackLogger)) // writer
	if vars.Debug {
		writers = append(writers, zapcore.AddSync(os.Stdout))
	}
	multiWriter := zapcore.NewMultiWriteSyncer(writers...)

	core := zapcore.NewCore(encoder, multiWriter, atomicLevel)

	var options []zap.Option
	options = append(options, zap.AddCaller())      // stack trace: show file name and line number
	options = append(options, zap.AddCallerSkip(1)) // stack trace: skip a layer because of package log's call
	if vars.Debug {
		options = append(options, zap.Development()) // development mode: makes DPanic-level logs panic instead of simply logging an error
	}
	options = append(options, zap.Fields(zap.String(serverKey, vars.ServerName))) // set filer "server" : ServerName

	zapLogger := zap.New(core, options...)
	logger = zapLogger.Sugar()
}

func Debug(args ...any) {
	logger.Debug(args...)
}

func Debugf(template string, args ...any) {
	logger.Debugf(template, args...)
}

func Debugw(msg string, keysAndValues ...any) {
	logger.Debugw(msg, keysAndValues...)
}

func Info(args ...any) {
	logger.Info(args...)
}

func Infof(template string, args ...any) {
	logger.Infof(template, args...)
}

func Infow(msg string, keysAndValues ...any) {
	logger.Infow(msg, keysAndValues...)
}

func Warn(args ...any) {
	logger.Warn(args...)
}

func Warnf(template string, args ...any) {
	logger.Warnf(template, args...)
}

func Warnw(msg string, keysAndValues ...any) {
	logger.Warnw(msg, keysAndValues...)
}

func Error(args ...any) {
	logger.Error(args...)
}

func Errorf(template string, args ...any) {
	logger.Errorf(template, args...)
}

func Errorw(msg string, keysAndValues ...any) {
	logger.Errorw(msg, keysAndValues...)
}

func DPanic(args ...any) {
	logger.DPanic(args...)
}

func DPanicf(template string, args ...any) {
	logger.DPanicf(template, args...)
}

func DPanicw(msg string, keysAndValues ...any) {
	logger.DPanicw(msg, keysAndValues...)
}

func Panic(args ...any) {
	logger.Panic(args...)
}

func Panicf(template string, args ...any) {
	logger.Panicf(template, args...)
}

func Panicw(msg string, keysAndValues ...any) {
	logger.Panicw(msg, keysAndValues...)
}

func Fatal(args ...any) {
	logger.Fatal(args...)
}

func Fatalf(template string, args ...any) {
	logger.Fatalf(template, args...)
}

func Fatalw(msg string, keysAndValues ...any) {
	logger.Fatalw(msg, keysAndValues...)
}

func Sync() {
	_ = logger.Sync()
}
