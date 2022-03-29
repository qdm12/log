package caller

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mutexBoolDifferentAddresses(t *testing.T, boolA, boolB *bool) {
	t.Helper()

	if boolA == nil || boolB == nil {
		return
	}

	addressA := fmt.Sprintf("%p", boolA)
	addressB := fmt.Sprintf("%p", boolB)
	assert.NotEqual(t, addressA, addressB)
}

func Test_Settings_SetDefaults(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initialSettings  Settings
		expectedSettings Settings
	}{
		"empty settings": {
			expectedSettings: Settings{
				File: boolPtr(false),
				Line: boolPtr(false),
				Func: boolPtr(false),
			},
		},
		"filled settings": {
			initialSettings: Settings{
				File: boolPtr(true),
				Line: boolPtr(true),
				Func: boolPtr(false),
			},
			expectedSettings: Settings{
				File: boolPtr(true),
				Line: boolPtr(true),
				Func: boolPtr(false),
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			settings := testCase.initialSettings
			settings.SetDefaults()

			assert.Equal(t, testCase.expectedSettings, settings)
		})
	}
}

func Test_Settings_Copy(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initialSettings  Settings
		expectedSettings Settings
	}{
		"empty settings": {},
		"filled settings": {
			initialSettings: Settings{
				File: boolPtr(true),
				Line: boolPtr(true),
				Func: boolPtr(false),
			},
			expectedSettings: Settings{
				File: boolPtr(true),
				Line: boolPtr(true),
				Func: boolPtr(false),
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			settings := testCase.initialSettings
			settingsCopy := settings.Copy()

			assert.Equal(t, testCase.expectedSettings, settingsCopy)

			mutexBoolDifferentAddresses(t, settings.File, settingsCopy.File)
			mutexBoolDifferentAddresses(t, settings.Line, settingsCopy.Line)
			mutexBoolDifferentAddresses(t, settings.Func, settingsCopy.Func)
		})
	}
}

func Test_Settings_OverrideWith(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		initialSettings  Settings
		otherSettings    Settings
		expectedSettings Settings
	}{
		"empty settings and empty override": {},
		"full settings and empty override": {
			initialSettings: Settings{
				File: boolPtr(true),
				Line: boolPtr(true),
				Func: boolPtr(false),
			},
			expectedSettings: Settings{
				File: boolPtr(true),
				Line: boolPtr(true),
				Func: boolPtr(false),
			},
		},
		"empty settings and full override": {
			otherSettings: Settings{
				File: boolPtr(true),
				Line: boolPtr(true),
				Func: boolPtr(false),
			},
			expectedSettings: Settings{
				File: boolPtr(true),
				Line: boolPtr(true),
				Func: boolPtr(false),
			},
		},
		"full settings and full override": {
			initialSettings: Settings{
				File: boolPtr(true),
				Line: boolPtr(true),
				Func: boolPtr(false),
			},
			otherSettings: Settings{
				File: boolPtr(false),
				Line: boolPtr(false),
				Func: boolPtr(true),
			},
			expectedSettings: Settings{
				File: boolPtr(false),
				Line: boolPtr(false),
				Func: boolPtr(true),
			},
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			settings := testCase.initialSettings
			settings.OverrideWith(testCase.otherSettings)

			assert.Equal(t, testCase.expectedSettings, settings)

			// Make sure we deep copy the other settings pointers
			mutexBoolDifferentAddresses(t, settings.File, testCase.otherSettings.File)
			mutexBoolDifferentAddresses(t, settings.Line, testCase.otherSettings.Line)
			mutexBoolDifferentAddresses(t, settings.Func, testCase.otherSettings.Func)
		})
	}
}

func Test_Line(t *testing.T) {
	t.Parallel()

	// find line number of log call below
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	b, err := os.ReadFile(file)
	require.NoError(t, err)
	first := true
	lineNumber := 0
	for i, line := range strings.Split(string(b), "\n") {
		if strings.Contains(line, "wrapFunc1()") {
			if first { // ignore the line above
				first = false
				continue
			}
			lineNumber = i + 1
			break
		}
	}

	testCases := map[string]struct {
		settings   Settings
		callerLine string
	}{
		"no show": {
			settings: Settings{
				File: boolPtr(false),
				Line: boolPtr(false),
				Func: boolPtr(false),
			},
		},
		"show file line": {
			settings: Settings{
				File: boolPtr(true),
				Line: boolPtr(true),
				Func: boolPtr(false),
			},
			callerLine: fmt.Sprintf("caller_test.go:L%d", lineNumber),
		},
		"show all": {
			settings: Settings{
				File: boolPtr(true),
				Line: boolPtr(true),
				Func: boolPtr(true),
			},
			callerLine: fmt.Sprintf("caller_test.go:L%d:func1", lineNumber),
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var callerLine string
			wrapFunc1 := func() { // Debug/Info calls
				func() { // log function
					callerLine = Line(testCase.settings)
				}()
			}

			wrapFunc1()

			assert.Equal(t, testCase.callerLine, callerLine)
		})
	}
}
