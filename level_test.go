package log

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Level_String(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		level Level
		s     string
	}{
		"debug": {
			level: LevelDebug,
			s:     "DEBUG",
		},
		"info": {
			level: LevelInfo,
			s:     "INFO",
		},
		"warn": {
			level: LevelWarn,
			s:     "WARN",
		},
		"error": {
			level: LevelError,
			s:     "ERROR",
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := testCase.level.String()

			assert.Equal(t, testCase.s, s)
		})
	}

	t.Run("invalid level", func(t *testing.T) {
		t.Parallel()

		invalidLevel := Level(99)

		assert.PanicsWithValue(t, "level 99 is unknown", func() {
			_ = invalidLevel.String()
		})
	})
}

func Test_Level_ColoredString(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		level Level
		s     string
	}{
		"debug": {
			level: LevelDebug,
			s:     "DEBUG",
		},
		"info": {
			level: LevelInfo,
			s:     "INFO",
		},
		"warn": {
			level: LevelWarn,
			s:     "WARN",
		},
		"error": {
			level: LevelError,
			s:     "ERROR",
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := testCase.level.ColoredString()
			// Note: fatih/color is clever enough to not add colors
			// when running tests, so the string is effectively without
			// color here.

			assert.Equal(t, testCase.s, s)
		})
	}
}

func Test_ParseLevel(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		s     string
		level Level
		err   error
	}{
		"debug": {
			s:     "DEBUG",
			level: LevelDebug,
		},
		"info": {
			s:     "INFO",
			level: LevelInfo,
		},
		"warn": {
			s:     "WARN",
			level: LevelWarn,
		},
		"error": {
			s:     "ERROR",
			level: LevelError,
		},
		"invalid": {
			s:   "someinvalid",
			err: errors.New("level is not recognized: someinvalid"),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			level, err := ParseLevel(testCase.s)

			if testCase.err != nil {
				require.EqualError(t, err, testCase.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.level, level)
		})
	}
}
