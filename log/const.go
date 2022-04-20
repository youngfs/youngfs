package log

const (
	timeFormat   = "2006-01-02 15:04:05.000000" //time.RFC3339Nano
	logExtension = ".log"

	timeKey       = "time"
	levelKey      = "level"
	nameKey       = "logger"
	callerKey     = "caller"
	messageKey    = "msg"
	stacktraceKey = "stacktrace"
	serverKey     = "server"
)
