package log

import "github.com/qdm12/log/internal/caller"

// RFC3339 format.
const timePrefixRegex = `^2[0-9]{3}-[0-1][0-9]-[0-3][0-9]T[0-2][0-9]:[0-5][0-9]:[0-5][0-9]Z `

func boolPtr(b bool) *bool { return &b }

func levelPtr(l Level) *Level { return &l }

func stringPtr(s string) *string { return &s }

func newCallerSettings(file, line, funC bool) caller.Settings {
	return caller.Settings{
		File: &file,
		Line: &line,
		Func: &funC,
	}
}
