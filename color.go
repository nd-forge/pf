package pf

import "strings"

// ANSI color codes
const (
	cReset   = "\033[0m"
	cKey     = "\033[36m" // cyan - field names
	cString  = "\033[32m" // green - strings
	cNumber  = "\033[33m" // yellow - numbers
	cBool    = "\033[35m" // magenta - booleans
	cType    = "\033[90m" // gray - type annotations
	cNil     = "\033[31m" // red - nil values
	cBrace   = "\033[37m" // white - braces/brackets
	cDiffDel = "\033[31m" // red - diff deletions
	cDiffAdd = "\033[32m" // green - diff additions
)

func (f *formatter) colored(color, text string) {
	if f.config.ColorMode {
		f.sb.WriteString(color)
		f.sb.WriteString(text)
		f.sb.WriteString(cReset)
	} else {
		f.sb.WriteString(text)
	}
}

func coloredStr(color, text string, colorMode bool) string {
	if colorMode {
		var sb strings.Builder
		sb.WriteString(color)
		sb.WriteString(text)
		sb.WriteString(cReset)
		return sb.String()
	}
	return text
}
