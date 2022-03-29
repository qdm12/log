package log

import (
	"fmt"
	"io"
	"time"

	"github.com/fatih/color"
	"github.com/qdm12/log/internal/caller"
)

func (l *Logger) logf(logLevel Level, format string, args ...interface{}) {
	l.settingsMutex.RLock()
	defer l.settingsMutex.RUnlock()
	settings := l.settings.copy()

	if *settings.level < logLevel {
		return
	}

	var line string
	if *settings.timeFormat != "" {
		line += time.Now().Format(*settings.timeFormat) + " "
	}

	line += logLevel.ColoredString() + " "
	if settings.component != "" {
		line += "[" + settings.component + "] "
	}

	if len(args) == 0 {
		line += format
	} else {
		line += fmt.Sprintf(format, args...)
	}

	callerString := caller.Line(settings.caller)
	if callerString != "" {
		line += "\t" + color.HiWhiteString(callerString)
	}

	line += "\n"

	l.writersMutexesMutex.RLock()
	for i, writer := range settings.writers {
		writerMutex := l.writersMutexes[i]
		if writerMutex == nil {
			// no need for a mutex, for example with io.Discard
			_, _ = io.WriteString(writer, line)
		} else {
			writerMutex.Lock()
			_, _ = io.WriteString(writer, line)
			writerMutex.Unlock()
		}
	}
	l.writersMutexesMutex.RUnlock()
}

// Debug logs with the debug level.
func (l *Logger) Debug(s string) { l.logf(LevelDebug, s) }

// Info logs with the info level.
func (l *Logger) Info(s string) { l.logf(LevelInfo, s) }

// Warn logs with the warn level.
func (l *Logger) Warn(s string) { l.logf(LevelWarn, s) }

// Error logs with the error level.
func (l *Logger) Error(s string) { l.logf(LevelError, s) }

// Debugf formats and logs at the debug level.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.logf(LevelDebug, format, args...)
}

// Infof formats and logs at the info level.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.logf(LevelInfo, format, args...)
}

// Warnf formats and logs at the warn level.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.logf(LevelWarn, format, args...)
}

// Errorf formats and logs at the error level.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.logf(LevelError, format, args...)
}
