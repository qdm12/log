package log

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// Level is the level of the logger.
type Level uint8

const (
	// LevelError is the error level.
	LevelError Level = iota
	// LevelWarn is the warn level.
	LevelWarn
	// LevelInfo is the info level.
	LevelInfo
	// LevelDebug is the debug level.
	LevelDebug
)

func (level Level) String() (s string) {
	switch level {
	case LevelError:
		return "ERROR"
	case LevelWarn:
		return "WARN"
	case LevelInfo:
		return "INFO"
	case LevelDebug:
		return "DEBUG"
	default:
		panic(fmt.Sprintf("level %d is unknown", level))
	}
}

// ColoredString returns the corresponding colored
// string for the level.
func (level Level) ColoredString() (s string) {
	attribute := color.Reset

	switch level {
	case LevelDebug:
		attribute = color.FgHiBlue
	case LevelInfo:
		attribute = color.FgCyan
	case LevelWarn:
		attribute = color.FgYellow
	case LevelError:
		attribute = color.FgHiRed
	}

	c := color.New(attribute)
	return c.Sprint(level.String())
}

var (
	ErrLevelNotRecognized = errors.New("level is not recognized")
)

// ParseLevel parses a string into a level, and returns an
// error if it fails.
func ParseLevel(s string) (level Level, err error) {
	switch strings.ToUpper(s) {
	case LevelDebug.String():
		return LevelDebug, nil
	case LevelInfo.String():
		return LevelInfo, nil
	case LevelWarn.String():
		return LevelWarn, nil
	case LevelError.String():
		return LevelError, nil
	}
	return 0, fmt.Errorf("%w: %s", ErrLevelNotRecognized, s)
}
