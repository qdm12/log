package log

import (
	"io"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Logger_Patch(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initialLogger  *Logger
		options        []Option
		expectedLogger *Logger
	}{
		"without option": {
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
		"with options": {
			initialLogger: &Logger{
				settings: settings{
					writers: []io.Writer{io.Discard},
					level:   levelPtr(LevelInfo),
					caller:  newCallerSettings(false, false, false),
				},
				writersMutexes: []*sync.Mutex{nil},
			},
			options: []Option{
				SetWriters(os.Stdout),
				SetLevel(LevelWarn),
				SetCallerFile(true),
			},
			expectedLogger: &Logger{
				settings: settings{
					writers: []io.Writer{os.Stdout},
					level:   levelPtr(LevelWarn),
					caller:  newCallerSettings(true, false, false),
				},
				writersMutexes: []*sync.Mutex{new(sync.Mutex)},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger := testCase.initialLogger

			logger.Patch(testCase.options...)

			assert.Equal(t,
				testCase.expectedLogger.settings,
				logger.settings)
		})
	}
}

func Test_Logger_patch(t *testing.T) {
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
		"with options": {
			initialLogger: &Logger{
				settings: settings{
					writers: []io.Writer{io.Discard},
					level:   levelPtr(LevelInfo),
					caller:  newCallerSettings(false, false, false),
				},
				writersMutexes: []*sync.Mutex{nil},
			},
			options: []Option{
				SetLevel(LevelWarn),
				SetCallerFile(true),
			},
			expectedLogger: &Logger{
				settings: settings{
					writers: []io.Writer{io.Discard},
					level:   levelPtr(LevelWarn),
					caller:  newCallerSettings(true, false, false),
				},
				writersMutexes: []*sync.Mutex{nil},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			logger := testCase.initialLogger

			logger.Patch(testCase.options...)

			assert.Equal(t, testCase.expectedLogger, logger)
		})
	}
}
