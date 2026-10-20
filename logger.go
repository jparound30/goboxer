package goboxer

type Logger interface {
	RequestDumpf(format string, args ...any)
	ResponseDumpf(format string, args ...any)
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
	EnabledLoggingResponseBody() bool
	EnabledLoggingRequestBody() bool
}
