package log

type Logger interface {
	Debug(args ...any)
	Debugf(template string, args ...any)
	Debugw(msg string, keysAndValues ...any)
	Info(args ...any)
	Infof(template string, args ...any)
	Infow(msg string, keysAndValues ...any)
	Warn(args ...any)
	Warnf(template string, args ...any)
	Warnw(msg string, keysAndValues ...any)
	Error(args ...any)
	Errorf(template string, args ...any)
	Errorw(msg string, keysAndValues ...any)
	DPanic(args ...any)
	DPanicf(template string, args ...any)
	DPanicw(msg string, keysAndValues ...any)
	Panic(args ...any)
	Panicf(template string, args ...any)
	Panicw(msg string, keysAndValues ...any)
	Fatal(args ...any)
	Fatalf(template string, args ...any)
	Fatalw(msg string, keysAndValues ...any)
	Sync() error
}
