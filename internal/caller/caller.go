package caller

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

type Settings struct {
	File *bool
	Line *bool
	Func *bool
}

func (s *Settings) SetDefaults() {
	s.File = defaultBoolPtr(s.File, false)
	s.Line = defaultBoolPtr(s.Line, false)
	s.Func = defaultBoolPtr(s.Func, false)
}

func (s *Settings) Copy() (settingsCopy Settings) {
	settingsCopy.File = copyBoolPtr(s.File)
	settingsCopy.Line = copyBoolPtr(s.Line)
	settingsCopy.Func = copyBoolPtr(s.Func)
	return settingsCopy
}

func (s *Settings) OverrideWith(other Settings) {
	s.File = overrideBoolPtr(s.File, other.File)
	s.Line = overrideBoolPtr(s.Line, other.Line)
	s.Func = overrideBoolPtr(s.Func, other.Func)
}

func Line(settings Settings) (s string) {
	if !*settings.File && !*settings.Line && !*settings.Func {
		return ""
	}

	const depth = 3
	pc, file, line, ok := runtime.Caller(depth)
	if !ok {
		return "error"
	}

	var fields []string

	if *settings.File {
		fields = append(fields, filepath.Base(file))
	}

	if *settings.Line {
		fields = append(fields, "L"+fmt.Sprint(line))
	}

	if *settings.Func {
		details := runtime.FuncForPC(pc)
		if details != nil {
			funcName := strings.TrimLeft(filepath.Ext(details.Name()), ".")
			fields = append(fields, funcName)
		}
	}

	return strings.Join(fields, ":")
}
