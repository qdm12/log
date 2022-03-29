package log

import (
	"io"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		options        []Option
		expectedLogger *Logger
	}{
		"no option": {
			expectedLogger: &Logger{
				settings: settings{
					writers:    []io.Writer{os.Stdout},
					level:      levelPtr(LevelInfo),
					timeFormat: stringPtr(time.RFC3339),
					caller:     newCallerSettings(false, false, false),
				},
				writersMutexes: []*sync.Mutex{new(sync.Mutex)},
			},
		},
		"all options": {
			options: []Option{
				SetLevel(LevelInfo),
				SetCallerFile(true),
				SetCallerLine(true),
				SetCallerFunc(true),
				SetTimeFormat(time.RFC1123),
				SetWriters(io.Discard),
			},
			expectedLogger: &Logger{
				settings: settings{
					writers:    []io.Writer{io.Discard},
					level:      levelPtr(LevelInfo),
					timeFormat: stringPtr(time.RFC1123),
					caller:     newCallerSettings(true, true, true),
				},
				writersMutexes: []*sync.Mutex{nil},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger := New(testCase.options...)

			assert.Equal(t, testCase.expectedLogger, logger)
		})
	}
}

func Test_Logger_New(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initialLogger  *Logger
		options        []Option
		expectedLogger *Logger
	}{
		"no option": {
			initialLogger: &Logger{
				settings: settings{
					writers: []io.Writer{os.Stdout},
					level:   levelPtr(LevelInfo),
					caller:  newCallerSettings(false, false, false),
				},
				writersMutexes: []*sync.Mutex{new(sync.Mutex)},
			},
			expectedLogger: &Logger{
				settings: settings{
					writers: []io.Writer{os.Stdout},
					level:   levelPtr(LevelInfo),
					caller:  newCallerSettings(false, false, false),
				},
				writersMutexes: []*sync.Mutex{new(sync.Mutex)},
			},
		},
		"some options": {
			initialLogger: &Logger{
				settings: settings{
					writers: []io.Writer{os.Stdout},
					level:   levelPtr(LevelInfo),
					caller:  newCallerSettings(true, true, true),
				},
				writersMutexes: []*sync.Mutex{new(sync.Mutex)},
			},
			options: []Option{
				SetLevel(LevelInfo),
				SetCallerFunc(false),
				SetTimeFormat(time.RFC1123),
				SetWriters(os.Stderr),
			},
			expectedLogger: &Logger{
				settings: settings{
					writers:    []io.Writer{os.Stderr},
					level:      levelPtr(LevelInfo),
					timeFormat: stringPtr(time.RFC1123),
					caller:     newCallerSettings(true, true, false),
				},
				writersMutexes: []*sync.Mutex{new(sync.Mutex)},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger := testCase.initialLogger.New(testCase.options...)

			assert.Equal(t, testCase.expectedLogger, logger)
		})
	}
}
