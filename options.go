package log

import (
	"io"
)

// Option is the type to specify settings modifier
// for the logger operation.
type Option func(s *settings)

// SetLevel sets the level for the logger.
// The level defaults to the lowest level, error.
func SetLevel(level Level) Option {
	return func(s *settings) {
		s.level = &level
	}
}

// SetComponent sets the component for the logger
// which will be logged on every log operation.
// Set it to the empty string so no component is logged.
// The default is the empty string component.
func SetComponent(component string) Option {
	return func(s *settings) {
		s.component = component
	}
}

// SetCallerFile enables or disables logging the caller file.
// The default is disabled.
func SetCallerFile(enabled bool) Option {
	return func(s *settings) {
		s.caller.File = &enabled
	}
}

// SetCallerLine enables or disables logging the caller line number.
// The default is disabled.
func SetCallerLine(enabled bool) Option {
	return func(s *settings) {
		s.caller.Line = &enabled
	}
}

// SetCallerFunc enables or disables logging the caller function.
// The default is disabled.
func SetCallerFunc(enabled bool) Option {
	return func(s *settings) {
		s.caller.Func = &enabled
	}
}

// SetTimeFormat set the time format for the logger.
// You can set it to an empty string in order to not
// log the time.
// The time format defaults to time.RFC3339.
func SetTimeFormat(timeFormat string) Option {
	return func(s *settings) {
		s.timeFormat = &timeFormat
	}
}

// SetWriters sets the writers for the logger.
// The writers defaults to a single writer of os.Stdout.
func SetWriters(writers ...io.Writer) Option {
	return func(s *settings) {
		s.writers = make([]io.Writer, len(writers))
		for i := range writers {
			s.writers[i] = writers[i]
		}
	}
}

// AddWriters adds the writers given to the existing
// writers for the logger.
// The writers defaults to a single writer of os.Stdout.
func AddWriters(writers ...io.Writer) Option {
	return func(s *settings) {
		newWriters := make([]io.Writer, len(s.writers), len(s.writers)+len(writers))
		copy(newWriters, s.writers)
		for _, writerToAdd := range writers {
			alreadyExists := false
			for _, existingWriter := range s.writers {
				if writerToAdd == existingWriter {
					alreadyExists = true
				}
			}
			if !alreadyExists {
				newWriters = append(newWriters, writerToAdd)
			}
		}
		s.writers = newWriters
	}
}
