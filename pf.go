// Package pf extends fmt with pretty-printing for Go values.
//
// Basic usage:
//
//	pf.Print(myStruct)
//	s := pf.Sprint(myStruct)
//	pf.Diff(oldStruct, newStruct)
//
// Implement PrettyPrinter for custom formatting:
//
//	func (u User) PrettyPrint() string { return fmt.Sprintf("User<%s>", u.Name) }
package pf

import (
	"fmt"
	"io"
	"os"
)

// DefaultConfig is the default pretty-print configuration.
var DefaultConfig = Config{
	Indent:      "  ",
	ShowTypes:   false,
	UseJSONTags: false,
	MaxDepth:    0,
	ColorMode:   true,
}

// --- Pretty Print ---

// Print pretty-prints to stdout.
func Print(v interface{}) {
	fmt.Fprintln(os.Stdout, Sprint(v))
}

// Sprint returns a pretty-printed string.
func Sprint(v interface{}) string {
	return DefaultConfig.Sprint(v)
}

// Fprint pretty-prints to the given writer.
func Fprint(w io.Writer, v interface{}) {
	fmt.Fprintln(w, DefaultConfig.Sprint(v))
}

// --- Diff ---

// Diff prints a colorized diff to stdout.
func Diff(a, b interface{}) {
	fmt.Fprintln(os.Stdout, SprintDiff(a, b))
}

// SprintDiff returns a diff string.
func SprintDiff(a, b interface{}) string {
	return DefaultConfig.SprintDiff(a, b)
}

// FprintDiff writes a diff to the given writer.
func FprintDiff(w io.Writer, a, b interface{}) {
	fmt.Fprintln(w, DefaultConfig.SprintDiff(a, b))
}
