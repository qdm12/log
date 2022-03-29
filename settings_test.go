package log

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/qdm12/log/internal/caller"
	"github.com/stretchr/testify/assert"
)

func Test_newSettings(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		options  []Option
		settings settings
	}{
		"no option": {},
		"multiple options": {
			options: []Option{
				SetLevel(LevelInfo),
				SetWriters(os.Stderr, os.Stdout),
			},
			settings: settings{
				level:   levelPtr(LevelInfo),
				writers: []io.Writer{os.Stderr, os.Stdout},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			settings := newSettings(testCase.options)

			assert.Equal(t, testCase.settings, settings)
		})
	}
}

func Test_settings_setDefaults(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initialSettings  settings
		expectedSettings settings
	}{
		"empty settings": {
			expectedSettings: settings{
				writers:    []io.Writer{os.Stdout},
				level:      levelPtr(LevelInfo),
				timeFormat: stringPtr(time.RFC3339),
				caller: caller.Settings{
					File: boolPtr(false),
					Line: boolPtr(false),
					Func: boolPtr(false),
				},
			},
		},
		"filled settings": {
			initialSettings: settings{
				writers:    []io.Writer{io.Discard},
				level:      levelPtr(LevelWarn),
				timeFormat: stringPtr(time.RFC1123),
				caller: caller.Settings{
					File: boolPtr(true),
					Line: boolPtr(true),
					Func: boolPtr(true),
				},
			},
			expectedSettings: settings{
				writers:    []io.Writer{io.Discard},
				level:      levelPtr(LevelWarn),
				timeFormat: stringPtr(time.RFC1123),
				caller: caller.Settings{
					File: boolPtr(true),
					Line: boolPtr(true),
					Func: boolPtr(true),
				},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			settings := testCase.initialSettings
			settings.setDefaults()

			assert.Equal(t, testCase.expectedSettings, settings)
		})
	}
}

func Test_settings_copy(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initialSettings  settings
		expectedSettings settings
	}{
		"empty settings": {},
		"filled settings": {
			initialSettings: settings{
				writers:    []io.Writer{io.Discard},
				level:      levelPtr(LevelWarn),
				timeFormat: stringPtr(time.RFC1123),
				component:  "component",
				caller: caller.Settings{
					File: boolPtr(true),
					Line: boolPtr(true),
					Func: boolPtr(true),
				},
			},
			expectedSettings: settings{
				writers:    []io.Writer{io.Discard},
				level:      levelPtr(LevelWarn),
				timeFormat: stringPtr(time.RFC1123),
				component:  "component",
				caller: caller.Settings{
					File: boolPtr(true),
					Line: boolPtr(true),
					Func: boolPtr(true),
				},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			settingsCopy := testCase.initialSettings.copy()

			assert.Equal(t, testCase.expectedSettings, settingsCopy)
		})
	}
}

func Test_settings_overrideWith(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initialSettings  settings
		overrideSettings settings
		expectedSettings settings
	}{
		"empty settings override with empty settings": {},
		"empty settings override with full settings": {
			overrideSettings: settings{
				writers:    []io.Writer{io.Discard},
				level:      levelPtr(LevelWarn),
				timeFormat: stringPtr(time.RFC1123),
				component:  "component",
				caller: caller.Settings{
					File: boolPtr(true),
					Line: boolPtr(true),
					Func: boolPtr(true),
				},
			},
			expectedSettings: settings{
				writers:    []io.Writer{io.Discard},
				level:      levelPtr(LevelWarn),
				timeFormat: stringPtr(time.RFC1123),
				component:  "component",
				caller: caller.Settings{
					File: boolPtr(true),
					Line: boolPtr(true),
					Func: boolPtr(true),
				},
			},
		},
		"filled settings": {
			initialSettings: settings{
				writers:    []io.Writer{io.Discard},
				level:      levelPtr(LevelWarn),
				timeFormat: stringPtr(time.RFC1123),
				component:  "component",
				caller: caller.Settings{
					File: boolPtr(false),
					Line: boolPtr(false),
					Func: boolPtr(false),
				},
			},
			overrideSettings: settings{
				writers:    []io.Writer{os.Stdout},
				level:      levelPtr(LevelInfo),
				timeFormat: stringPtr(time.RFC3339),
				component:  "new component",
				caller: caller.Settings{
					File: boolPtr(true),
					Line: boolPtr(true),
					Func: boolPtr(true),
				},
			},
			expectedSettings: settings{
				writers:    []io.Writer{os.Stdout},
				level:      levelPtr(LevelInfo),
				timeFormat: stringPtr(time.RFC3339),
				component:  "new component",
				caller: caller.Settings{
					File: boolPtr(true),
					Line: boolPtr(true),
					Func: boolPtr(true),
				},
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			settings := testCase.initialSettings
			settings.overrideWith(testCase.overrideSettings)

			assert.Equal(t, testCase.expectedSettings, settings)
		})
	}
}
