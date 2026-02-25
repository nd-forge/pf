package main

import (
	"fmt"

	"github.com/taiki-nd/pf"
)

// --- Custom interface examples ---

type Token struct {
	Value string
}

// PrettyPrinter: mask sensitive data
func (t Token) PrettyPrint() string {
	if len(t.Value) > 4 {
		return fmt.Sprintf("Token(***%s)", t.Value[len(t.Value)-4:])
	}
	return "Token(****)"
}

type Status int

// fmt.Stringer
func (s Status) String() string {
	names := map[Status]string{0: "Inactive", 1: "Active", 2: "Suspended"}
	if n, ok := names[s]; ok {
		return n
	}
	return "Unknown"
}

// --- Data types ---

type Address struct {
	City    string `json:"city"`
	Country string `json:"country"`
	Zip     string `json:"zip,omitempty"`
}

type Order struct {
	ID     int     `json:"id"`
	Item   string  `json:"item"`
	Amount float64 `json:"amount"`
}

type User struct {
	Name    string            `json:"name"`
	Age     int               `json:"age"`
	Email   string            `json:"email,omitempty"`
	Status  Status            `json:"status"`
	Token   Token             `json:"-"`
	Address Address           `json:"address"`
	Tags    []string          `json:"tags"`
	Orders  []Order           `json:"orders"`
	Meta    map[string]string `json:"meta,omitempty"`
}

func main() {
	user := User{
		Name:   "John Smith",
		Age:    30,
		Status: 1,
		Token:  Token{Value: "secret_abcd1234"},
		Address: Address{
			City:    "San Francisco",
			Country: "USA",
			Zip:     "94105",
		},
		Tags: []string{"developer", "trader", "audiophile"},
		Orders: []Order{
			{ID: 1, Item: "McIntosh MA12000", Amount: 1650000},
			{ID: 2, Item: "B&W 802 D4", Amount: 3300000},
		},
		Meta: map[string]string{
			"role":   "admin",
			"region": "APAC",
		},
	}

	// ===== Basic pretty print =====
	fmt.Println("=== Print (default) ===")
	pf.Print(user)

	// ===== With type annotations =====
	fmt.Println("\n=== ShowTypes ===")
	cfg := pf.Config{
		Indent:    "  ",
		ShowTypes: true,
		ColorMode: true,
	}
	cfg.Print(user)

	// ===== JSON tags =====
	fmt.Println("\n=== UseJSONTags ===")
	jsonCfg := pf.Config{
		Indent:      "  ",
		UseJSONTags: true,
		ColorMode:   true,
	}
	jsonCfg.Print(user)

	// ===== No color (for logging) =====
	fmt.Println("\n=== No Color ===")
	logCfg := pf.Config{
		Indent:    "  ",
		ColorMode: false,
	}
	fmt.Println(logCfg.Sprint(user))

	// ===== Diff =====
	fmt.Println("\n=== Diff ===")
	oldUser := user
	newUser := user
	newUser.Age = 31
	newUser.Status = 2
	newUser.Tags = []string{"developer", "trader", "audiophile", "gopher"}
	newUser.Address.City = "New York"

	pf.Diff(oldUser, newUser)

	// ===== Interface demo: Stringer =====
	fmt.Println("\n=== Stringer (Status) ===")
	pf.Print(Status(1))

	// ===== Interface demo: PrettyPrinter =====
	fmt.Println("\n=== PrettyPrinter (Token) ===")
	pf.Print(Token{Value: "secret_abcd1234"})

	// ===== Sprint for logging =====
	fmt.Println("\n=== Sprint ===")
	pf.DefaultConfig.ColorMode = false
	s := pf.Sprint(map[string]interface{}{
		"action": "login",
		"user":   "jsmith",
		"ok":     true,
	})
	fmt.Printf("[LOG] request: %s\n", s)
}
