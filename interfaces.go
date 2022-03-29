package log

var _ LoggerInterface = (*Logger)(nil)

type LoggerInterface interface {
	LeveledLogger
	LoggerPatcher
	ChildConstructor
}

// LeveledLogger is the interface to log at different levels.
type LeveledLogger interface {
	Debug(s string)
	Info(s string)
	Warn(s string)
	Error(s string)
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// LoggerPatcher is the interface to update the current logger.
type LoggerPatcher interface {
	Patch(options ...Option)
}

// ChildConstructor is the interface to create child loggers.
type ChildConstructor interface {
	New(options ...Option) *Logger
}
