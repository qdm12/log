package log

import (
	"sync"

	"github.com/qdm12/log/internal/writersreg"
)

var writersRegistry = writersreg.NewRegistry() //nolint:gochecknoglobals

// Logger is the logger implementation structure.
// It is thread safe to use.
type Logger struct {
	settings      settings
	settingsMutex sync.RWMutex
	// writersMutexes is a slice of mutex pointers
	// matching the order of settings.writers.
	writersMutexes      []*sync.Mutex
	writersMutexesMutex sync.RWMutex
}

// New creates a new logger, with thread safety each of
// its writers and other loggers. You can pass options
// to configure the logger.
func New(options ...Option) *Logger {
	settings := newSettings(options)
	settings.setDefaults()

	writerMutexes := writersRegistry.RegisterWriters(settings.writers)

	return &Logger{
		settings:       settings,
		writersMutexes: writerMutexes,
	}
}

// New creates a child logger inheriting from the settings of
// the current logger. Options can be passed to modify
// the settings of the new child logger to be created.
func (l *Logger) New(options ...Option) *Logger {
	newSettings := newSettings(options)

	l.settingsMutex.RLock()
	childSettings := l.settings.copy()
	l.settingsMutex.RUnlock()

	childSettings.overrideWith(newSettings)
	// defaults are already set in parent

	writersMutexes := writersRegistry.RegisterWriters(childSettings.writers)

	return &Logger{
		settings:       childSettings,
		writersMutexes: writersMutexes,
	}
}
