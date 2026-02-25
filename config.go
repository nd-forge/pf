package pf

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

// Config controls pretty-print formatting.
type Config struct {
	// Indent string per level. Default: "  "
	Indent string
	// ShowTypes annotates values with their type name.
	ShowTypes bool
	// UseJSONTags uses json tag names instead of Go field names.
	UseJSONTags bool
	// MaxDepth limits nesting depth (0 = unlimited).
	MaxDepth int
	// ColorMode enables ANSI color output.
	ColorMode bool
}

// Sprint returns a pretty-printed string using this config.
func (c Config) Sprint(v interface{}) string {
	f := &formatter{config: c}
	f.format(reflect.ValueOf(v), 0)
	return f.sb.String()
}

// Print pretty-prints to stdout using this config.
func (c Config) Print(v interface{}) {
	fmt.Fprintln(os.Stdout, c.Sprint(v))
}

// Fprint pretty-prints to the given writer using this config.
func (c Config) Fprint(w io.Writer, v interface{}) {
	fmt.Fprintln(w, c.Sprint(v))
}

// SprintDiff returns a diff string using this config.
func (c Config) SprintDiff(a, b interface{}) string {
	d := &differ{config: c}
	return d.diff(a, b)
}

type formatter struct {
	config Config
	sb     strings.Builder
}
