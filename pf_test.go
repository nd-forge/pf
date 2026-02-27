package pf

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// --- Test types ---

type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

type User struct {
	Name    string  `json:"name"`
	Age     int     `json:"age"`
	Email   string  `json:"email,omitempty"`
	Active  bool    `json:"active"`
	Address Address `json:"address"`
	Tags    []string
}

// PrettyPrinter implementation
type Token struct {
	Value string
}

func (t Token) PrettyPrint() string {
	return fmt.Sprintf("Token(***%s)", t.Value[len(t.Value)-4:])
}

// Stringer implementation
type Status int

func (s Status) String() string {
	switch s {
	case 0:
		return "Inactive"
	case 1:
		return "Active"
	default:
		return "Unknown"
	}
}

// error implementation
type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// --- Tests ---

func TestPrint_BasicStruct(t *testing.T) {
	user := User{
		Name:   "John Smith",
		Age:    30,
		Active: true,
		Address: Address{
			City:    "San Francisco",
			Country: "USA",
		},
		Tags: []string{"go", "trading"},
	}

	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(user)

	expects := []string{
		`Name: "John Smith"`,
		`Age: 30`,
		`Active: true`,
		`City: "San Francisco"`,
		`Tags: ["go", "trading"]`,
	}
	for _, e := range expects {
		if !strings.Contains(got, e) {
			t.Errorf("expected %q in output, got:\n%s", e, got)
		}
	}
}

func TestPrint_ShowTypes(t *testing.T) {
	user := User{Name: "Test"}
	c := Config{Indent: "  ", ShowTypes: true, ColorMode: false}
	got := c.Sprint(user)

	if !strings.Contains(got, "User {") {
		t.Errorf("expected type annotation, got:\n%s", got)
	}
}

func TestPrint_JSONTags(t *testing.T) {
	user := User{
		Name:   "Test",
		Age:    25,
		Active: true,
	}
	c := Config{Indent: "  ", UseJSONTags: true, ColorMode: false}
	got := c.Sprint(user)

	if !strings.Contains(got, "name:") {
		t.Errorf("expected json tag 'name', got:\n%s", got)
	}
	// Email is omitempty and zero, should be omitted
	if strings.Contains(got, "email") {
		t.Errorf("expected email to be omitted (omitempty), got:\n%s", got)
	}
}

func TestPrint_Nil(t *testing.T) {
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(nil)
	if got != "nil" {
		t.Errorf("expected nil, got: %s", got)
	}
}

func TestPrint_NilPointer(t *testing.T) {
	var u *User
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(u)
	if got != "nil" {
		t.Errorf("expected nil, got: %s", got)
	}
}

func TestPrint_Map(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(m)

	if !strings.Contains(got, `"a": 1`) {
		t.Errorf("expected map entry, got:\n%s", got)
	}
}

func TestPrint_EmptySlice(t *testing.T) {
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint([]int{})
	if got != "[]" {
		t.Errorf("expected [], got: %s", got)
	}
}

func TestPrint_MaxDepth(t *testing.T) {
	user := User{
		Name:    "Test",
		Address: Address{City: "New York"},
	}
	c := Config{Indent: "  ", MaxDepth: 1, ColorMode: false}
	got := c.Sprint(user)

	if !strings.Contains(got, "...") {
		t.Errorf("expected ... for depth limit, got:\n%s", got)
	}
}

// --- Interface tests ---

func TestPrettyPrinter(t *testing.T) {
	tok := Token{Value: "abcd1234"}
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(tok)

	expected := "Token(***1234)"
	if got != expected {
		t.Errorf("expected %q, got: %q", expected, got)
	}
}

func TestStringer(t *testing.T) {
	s := Status(1)
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(s)

	if got != `"Active"` {
		t.Errorf("expected \"Active\", got: %q", got)
	}
}

func TestError(t *testing.T) {
	e := &AppError{Code: 404, Message: "not found"}
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(e)

	if !strings.Contains(got, "not found") {
		t.Errorf("expected error message, got: %q", got)
	}
}

// --- Diff tests ---

func TestDiff_Struct(t *testing.T) {
	a := User{Name: "John", Age: 30, Active: true}
	b := User{Name: "John", Age: 31, Active: false}

	c := Config{Indent: "  ", ColorMode: false}
	got := c.SprintDiff(a, b)

	expects := []string{
		`Name: "John"`, // unchanged
		"- Age: 30",
		"+ Age: 31",
		"- Active: true",
		"+ Active: false",
	}
	for _, e := range expects {
		if !strings.Contains(got, e) {
			t.Errorf("expected %q in diff, got:\n%s", e, got)
		}
	}
}

func TestDiff_Map(t *testing.T) {
	a := map[string]int{"x": 1, "y": 2}
	b := map[string]int{"x": 1, "z": 3}

	c := Config{Indent: "  ", ColorMode: false}
	got := c.SprintDiff(a, b)

	if !strings.Contains(got, "- y:") {
		t.Errorf("expected deleted key y, got:\n%s", got)
	}
	if !strings.Contains(got, "+ z:") {
		t.Errorf("expected added key z, got:\n%s", got)
	}
}

func TestDiff_Slice(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{1, 4, 3, 5}

	c := Config{Indent: "  ", ColorMode: false}
	got := c.SprintDiff(a, b)

	if !strings.Contains(got, "- [1]: 2") {
		t.Errorf("expected changed index, got:\n%s", got)
	}
	if !strings.Contains(got, "+ [1]: 4") {
		t.Errorf("expected changed index, got:\n%s", got)
	}
}

// --- Top-level pf.go function tests ---

func TestTopLevel_Sprint(t *testing.T) {
	got := Sprint(42)
	if !strings.Contains(got, "42") {
		t.Errorf("expected 42, got: %s", got)
	}
}

func TestTopLevel_Print(t *testing.T) {
	// Print writes to stdout; just verify it doesn't panic
	Print(User{Name: "Test"})
}

func TestTopLevel_Fprint(t *testing.T) {
	var buf bytes.Buffer
	Fprint(&buf, User{Name: "Test"})
	got := buf.String()
	if !strings.Contains(got, "Test") {
		t.Errorf("expected Test in output, got: %s", got)
	}
}

func TestTopLevel_SprintDiff(t *testing.T) {
	got := SprintDiff(1, 2)
	if !strings.Contains(got, "1") || !strings.Contains(got, "2") {
		t.Errorf("expected diff values, got: %s", got)
	}
}

func TestTopLevel_Diff(t *testing.T) {
	// Diff writes to stdout; just verify it doesn't panic
	Diff(1, 2)
}

func TestTopLevel_FprintDiff(t *testing.T) {
	var buf bytes.Buffer
	FprintDiff(&buf, 1, 2)
	got := buf.String()
	if !strings.Contains(got, "1") || !strings.Contains(got, "2") {
		t.Errorf("expected diff values, got: %s", got)
	}
}

// --- Config Print/Fprint tests ---

func TestConfig_Print(t *testing.T) {
	c := Config{Indent: "  ", ColorMode: false}
	// Just verify it doesn't panic
	c.Print(User{Name: "Test"})
}

func TestConfig_Fprint(t *testing.T) {
	c := Config{Indent: "  ", ColorMode: false}
	var buf bytes.Buffer
	c.Fprint(&buf, User{Name: "Test"})
	got := buf.String()
	if !strings.Contains(got, "Test") {
		t.Errorf("expected Test in output, got: %s", got)
	}
}

// --- Color mode tests ---

func TestPrint_ColorMode(t *testing.T) {
	user := User{
		Name:   "John",
		Age:    30,
		Active: true,
		Address: Address{
			City:    "SF",
			Country: "US",
		},
		Tags: []string{"go"},
	}
	c := Config{Indent: "  ", ColorMode: true}
	got := c.Sprint(user)

	// Should contain ANSI color codes
	if !strings.Contains(got, "\033[") {
		t.Errorf("expected ANSI color codes in output, got:\n%s", got)
	}
	// Should contain reset codes
	if !strings.Contains(got, cReset) {
		t.Errorf("expected reset code in output, got:\n%s", got)
	}
}

func TestPrint_ColorMode_Map(t *testing.T) {
	m := map[string]int{"a": 1}
	c := Config{Indent: "  ", ColorMode: true}
	got := c.Sprint(m)
	if !strings.Contains(got, "\033[") {
		t.Errorf("expected ANSI color codes in map output, got:\n%s", got)
	}
}

func TestPrint_ColorMode_Nil(t *testing.T) {
	c := Config{Indent: "  ", ColorMode: true}
	got := c.Sprint(nil)
	if !strings.Contains(got, cNil) {
		t.Errorf("expected nil color in output, got:\n%s", got)
	}
}

// --- Scalar type tests ---

func TestPrint_Float(t *testing.T) {
	c := Config{Indent: "  ", ColorMode: false}

	// Integer-valued float
	got := c.Sprint(42.0)
	if got != "42.0" {
		t.Errorf("expected 42.0, got: %s", got)
	}

	// Decimal float
	got = c.Sprint(3.14)
	if got != "3.14" {
		t.Errorf("expected 3.14, got: %s", got)
	}
}

func TestPrint_Uint(t *testing.T) {
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(uint(42))
	if got != "42" {
		t.Errorf("expected 42, got: %s", got)
	}
}

func TestPrint_Bool(t *testing.T) {
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(true)
	if got != "true" {
		t.Errorf("expected true, got: %s", got)
	}
}

func TestPrint_String(t *testing.T) {
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint("hello")
	if got != `"hello"` {
		t.Errorf("expected \"hello\", got: %s", got)
	}
}

// --- Chan and Func tests ---

func TestPrint_Chan(t *testing.T) {
	ch := make(chan int)
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(ch)
	if !strings.Contains(got, "chan") {
		t.Errorf("expected chan in output, got: %s", got)
	}
}

func TestPrint_Func(t *testing.T) {
	fn := func() {}
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(fn)
	if !strings.Contains(got, "func") {
		t.Errorf("expected func in output, got: %s", got)
	}
}

// --- Interface value test ---

func TestPrint_InterfaceValue(t *testing.T) {
	var iface interface{} = map[string]interface{}{
		"key": "value",
	}
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(iface)
	if !strings.Contains(got, "key") {
		t.Errorf("expected key in output, got: %s", got)
	}
}

func TestPrint_NilInterface(t *testing.T) {
	type Container struct {
		Value interface{}
	}
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(Container{Value: nil})
	if !strings.Contains(got, "nil") {
		t.Errorf("expected nil in output, got: %s", got)
	}
}

// --- Slice formatting tests ---

func TestPrint_LongSlice(t *testing.T) {
	// More than 5 simple elements => multi-line format
	s := []int{1, 2, 3, 4, 5, 6}
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(s)
	if !strings.Contains(got, "[\n") {
		t.Errorf("expected multi-line slice, got:\n%s", got)
	}
}

func TestPrint_StructSlice(t *testing.T) {
	// Slice of non-simple kinds => multi-line format
	s := []Address{{City: "SF"}, {City: "NY"}}
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(s)
	if !strings.Contains(got, "SF") || !strings.Contains(got, "NY") {
		t.Errorf("expected cities in output, got:\n%s", got)
	}
}

func TestPrint_NilSlice(t *testing.T) {
	var s []int
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(s)
	if got != "nil" {
		t.Errorf("expected nil, got: %s", got)
	}
}

func TestPrint_Array(t *testing.T) {
	a := [3]int{1, 2, 3}
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(a)
	if !strings.Contains(got, "1") && !strings.Contains(got, "2") && !strings.Contains(got, "3") {
		t.Errorf("expected array elements, got:\n%s", got)
	}
}

// --- Map edge cases ---

func TestPrint_NilMap(t *testing.T) {
	var m map[string]int
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(m)
	if got != "nil" {
		t.Errorf("expected nil, got: %s", got)
	}
}

func TestPrint_EmptyMap(t *testing.T) {
	m := map[string]int{}
	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(m)
	if got != "{}" {
		t.Errorf("expected {}, got: %s", got)
	}
}

func TestPrint_MapShowTypes(t *testing.T) {
	m := map[string]int{"a": 1}
	c := Config{Indent: "  ", ShowTypes: true, ColorMode: false}
	got := c.Sprint(m)
	if !strings.Contains(got, "map[string]int") {
		t.Errorf("expected type annotation, got:\n%s", got)
	}
}
