package log

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Logger_log(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		logger      *Logger
		level       Level
		s           string
		args        []interface{}
		outputRegex string
	}{
		"log at info with debug set": {
			logger: &Logger{
				settings: settings{
					writers:    []io.Writer{bytes.NewBuffer(nil)},
					level:      levelPtr(LevelDebug),
					timeFormat: stringPtr(time.RFC3339),
					caller:     newCallerSettings(false, false, false),
				},
				writersMutexes: []*sync.Mutex{new(sync.Mutex)},
			},
			level:       LevelInfo,
			s:           "some words",
			outputRegex: timePrefixRegex + "INFO some words\n$",
		},
		"log at info with warn set": {
			logger: &Logger{
				settings: settings{
					writers:    []io.Writer{bytes.NewBuffer(nil)},
					level:      levelPtr(LevelWarn),
					timeFormat: stringPtr(time.RFC3339),
					caller:     newCallerSettings(false, false, false),
				},
				writersMutexes: []*sync.Mutex{new(sync.Mutex)},
			},
			level:       LevelInfo,
			s:           "some words",
			outputRegex: "^$",
		},
		"format string": {
			logger: &Logger{
				settings: settings{
					writers:    []io.Writer{bytes.NewBuffer(nil)},
					level:      levelPtr(LevelDebug),
					timeFormat: stringPtr(time.RFC3339),
					caller:     newCallerSettings(false, false, false),
				},
				writersMutexes: []*sync.Mutex{new(sync.Mutex)},
			},
			level:       LevelDebug,
			s:           "some %s",
			args:        []interface{}{"words"},
			outputRegex: timePrefixRegex + "DEBUG some words\n$",
		},
		"show caller": {
			logger: &Logger{
				settings: settings{
					writers:    []io.Writer{bytes.NewBuffer(nil)},
					level:      levelPtr(LevelDebug),
					timeFormat: stringPtr(time.RFC3339),
					caller:     newCallerSettings(true, true, true),
				},
				writersMutexes: []*sync.Mutex{new(sync.Mutex)},
			},
			level:       LevelDebug,
			s:           "some words",
			outputRegex: timePrefixRegex + "DEBUG some words\tlog_test.go:L[0-9]+:func[0-9]+\n$",
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Retrieve buffer set in test case as writer
			require.NotEmpty(t, testCase.logger.settings.writers)
			buffer, ok := testCase.logger.settings.writers[0].(*bytes.Buffer)
			require.True(t, ok)

			logWrapper := func() { // wrap for caller depth of 3
				testCase.logger.logf(testCase.level, testCase.s, testCase.args...)
			}

			logWrapper()

			line := buffer.String()
			buffer.Reset()

			regex, err := regexp.Compile(testCase.outputRegex)
			require.NoError(t, err)

			assert.True(t, regex.MatchString(line),
				"line %q does not match regex %q", line, regex.String())
		})
	}
}

func Test_Logger_LevelsLog(t *testing.T) {
	t.Parallel()

	buffer := bytes.NewBuffer(nil)

	logger := New(SetLevel(LevelDebug), SetWriters(buffer))
	logger.Debug("some debug")
	logger.Info("some info")
	logger.Warn("some warn")
	logger.Error("some error")
	logger.Debugf("some %dnd debug", 2)
	logger.Infof("some %dnd info", 2)
	logger.Warnf("some %dnd warn", 2)
	logger.Errorf("some %dnd error", 2)

	lines := strings.Split(buffer.String(), "\n")
	buffer.Reset()

	// Check for trailing newline
	require.NotEmpty(t, lines)
	assert.Equal(t, "", lines[len(lines)-1])
	lines = lines[:len(lines)-1]

	expectedRegexes := []string{
		timePrefixRegex + "DEBUG some debug$",
		timePrefixRegex + "INFO some info$",
		timePrefixRegex + "WARN some warn$",
		timePrefixRegex + "ERROR some error$",
		timePrefixRegex + "DEBUG some 2nd debug$",
		timePrefixRegex + "INFO some 2nd info$",
		timePrefixRegex + "WARN some 2nd warn$",
		timePrefixRegex + "ERROR some 2nd error$",
	}

	require.Equal(t, len(expectedRegexes), len(lines))

	for i := range lines {
		regex, err := regexp.Compile(expectedRegexes[i])
		require.NoError(t, err)

		assert.True(t, regex.MatchString(lines[i]),
			"line %q does not match regex %q", lines[i], expectedRegexes[i])
	}
}
