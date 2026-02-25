package pf

import (
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
		Name:   "Noda",
		Age:    30,
		Active: true,
		Address: Address{
			City:    "Yokohama",
			Country: "Japan",
		},
		Tags: []string{"go", "trading"},
	}

	c := Config{Indent: "  ", ColorMode: false}
	got := c.Sprint(user)

	expects := []string{
		`Name: "Noda"`,
		`Age: 30`,
		`Active: true`,
		`City: "Yokohama"`,
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
		Address: Address{City: "Tokyo"},
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
	a := User{Name: "Noda", Age: 30, Active: true}
	b := User{Name: "Noda", Age: 31, Active: false}

	c := Config{Indent: "  ", ColorMode: false}
	got := c.SprintDiff(a, b)

	expects := []string{
		`Name: "Noda"`, // unchanged
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
