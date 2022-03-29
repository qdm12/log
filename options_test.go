package log

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/qdm12/log/internal/caller"
	"github.com/stretchr/testify/assert"
)

func Test_Option(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initialSettings  settings
		option           Option
		expectedSettings settings
	}{
		"SetLevel": {
			option: SetLevel(LevelInfo),
			expectedSettings: settings{
				level: levelPtr(LevelInfo),
			},
		},
		"SetCallerFile": {
			option: SetCallerFile(false),
			expectedSettings: settings{
				caller: caller.Settings{
					File: boolPtr(false),
				},
			},
		},
		"SetCallerLine": {
			option: SetCallerLine(true),
			expectedSettings: settings{
				caller: caller.Settings{
					Line: boolPtr(true),
				},
			},
		},
		"SetCallerFunc": {
			option: SetCallerFunc(true),
			expectedSettings: settings{
				caller: caller.Settings{
					Func: boolPtr(true),
				},
			},
		},
		"SetTimeFormat": {
			option: SetTimeFormat("123"),
			expectedSettings: settings{
				timeFormat: stringPtr("123"),
			},
		},
		"SetWriters": {
			option: SetWriters(os.Stdout, io.Discard),
			expectedSettings: settings{
				writers: []io.Writer{os.Stdout, io.Discard},
			},
		},
		"AddWriters": {
			initialSettings: settings{
				writers: []io.Writer{bytes.NewBuffer(nil), io.Discard},
			},
			option: AddWriters(os.Stdout, io.Discard),
			expectedSettings: settings{
				writers: []io.Writer{bytes.NewBuffer(nil), io.Discard, os.Stdout},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			settings := testCase.initialSettings
			testCase.option(&settings)

			assert.Equal(t, testCase.expectedSettings, settings)
		})
	}
}
