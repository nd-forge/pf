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
