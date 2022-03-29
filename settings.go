package log

import (
	"io"
	"os"
	"time"

	"github.com/qdm12/log/internal/caller"
)

type settings struct {
	writers    []io.Writer
	level      *Level
	timeFormat *string
	component  string
	caller     caller.Settings
}

// newSettings returns settings using the options given
// and with defaults set for any unset setting field.
func newSettings(options []Option) (s settings) {
	for _, option := range options {
		option(&s)
	}
	return s
}

func (s *settings) setDefaults() {
	if len(s.writers) == 0 {
		s.writers = []io.Writer{os.Stdout}
	}

	if s.level == nil {
		value := LevelInfo
		s.level = &value
	}

	if s.timeFormat == nil {
		value := time.RFC3339
		s.timeFormat = &value
	}

	s.caller.SetDefaults()
}

func (s *settings) copy() (settingsCopy settings) {
	if s.writers != nil {
		settingsCopy.writers = make([]io.Writer, len(s.writers))
		copy(settingsCopy.writers, s.writers)
	}

	if s.level != nil {
		level := *s.level
		settingsCopy.level = &level
	}

	if s.timeFormat != nil {
		timeFormat := *s.timeFormat
		settingsCopy.timeFormat = &timeFormat
	}

	settingsCopy.component = s.component

	settingsCopy.caller = s.caller.Copy()

	return settingsCopy
}

func (s *settings) overrideWith(other settings) {
	if len(other.writers) > 0 {
		s.writers = other.writers
	}

	if other.level != nil {
		value := *other.level
		s.level = &value
	}

	if other.timeFormat != nil {
		value := *other.timeFormat
		s.timeFormat = &value
	}

	if other.component != "" {
		s.component = other.component
	}

	s.caller.OverrideWith(other.caller)
}
