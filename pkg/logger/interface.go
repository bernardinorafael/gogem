package logger

type Logger interface {
	Debug(msg any, keyvals ...any)
	Debugf(format string, args ...any)
	Info(msg any, keyvals ...any)
	Infof(format string, args ...any)
	Warn(msg any, keyvals ...any)
	Warnf(format string, args ...any)
	Error(msg any, keyvals ...any)
	Errorf(format string, args ...any)
	Fatal(msg any, keyvals ...any)
	Fatalf(format string, args ...any)
	With(keyvals ...any) Logger
}
