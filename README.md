# pf

[![Go Reference](https://pkg.go.dev/badge/github.com/nd-forge/pf.svg)](https://pkg.go.dev/github.com/nd-forge/pf)
[![CI](https://github.com/nd-forge/pf/actions/workflows/ci.yml/badge.svg)](https://github.com/nd-forge/pf/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/nd-forge/pf/branch/main/graph/badge.svg)](https://codecov.io/gh/nd-forge/pf)
[![Go Report Card](https://goreportcard.com/badge/github.com/nd-forge/pf)](https://goreportcard.com/report/github.com/nd-forge/pf)

**English** | [日本語](README_ja.md) | [中文](README_zh.md) | [한국어](README_ko.md) | [Español](README_es.md) | [Português](README_pt.md)

A Go pretty-print package that extends `fmt`. Formats structs, maps, and slices with indentation for readable output.

## Install

```bash
go get github.com/nd-forge/pf
```

## Quick Start

```go
import "github.com/nd-forge/pf"

type User struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Active  bool
    Address Address
}

user := User{Name: "John", Age: 30, Active: true, Address: Address{City: "San Francisco"}}

// Pretty print
pf.Print(user)
```

Output (with ANSI colors in terminal):
```
{
  Name: "John",
  Age: 30,
  Active: true,
  Address: {
    City: "San Francisco"
  }
}
```

## API

### Pretty Print

| Function | Description |
|---|---|
| `pf.Print(v)` | Print to stdout |
| `pf.Sprint(v)` | Return as string |
| `pf.Fprint(w, v)` | Write to io.Writer |

### Diff

| Function | Description |
|---|---|
| `pf.Diff(a, b)` | Print diff to stdout |
| `pf.SprintDiff(a, b)` | Return diff as string |
| `pf.FprintDiff(w, a, b)` | Write diff to io.Writer |

```go
old := User{Name: "John", Age: 30, Active: true}
new := User{Name: "John", Age: 31, Active: false}

pf.Diff(old, new)
// {
//   Name: "John"
//   - Age: 30
//   + Age: 31
//   - Active: true
//   + Active: false
// }
```

## Config

```go
c := pf.Config{
    Indent:      "    ",  // 4-space indent
    ShowTypes:   true,    // show type names
    UseJSONTags: true,    // use `json:"..."` tag names
    MaxDepth:    3,       // limit nesting
    ColorMode:   false,   // no ANSI colors (for logging)
}

c.Print(myStruct)
```

### UseJSONTags

```go
type User struct {
    Name  string `json:"user_name"`
    Email string `json:"email_address,omitempty"`
}

// UseJSONTags: true
// {
//   user_name: "John",
// }
// (Email is omitempty + zero value → omitted)
```

## Interfaces

### PrettyPrinter

Highest priority. Define custom pretty-print output for your types.

```go
type Token struct { Value string }

func (t Token) PrettyPrint() string {
    return fmt.Sprintf("Token(***%s)", t.Value[len(t.Value)-4:])
}

pf.Print(Token{Value: "secret1234"})
// Token(***1234)
```

### PrettyPrinterConfig

Receives the Config to produce format-aware output.

```go
func (t Token) PrettyPrintConfig(c pf.Config) string {
    masked := "***" + t.Value[len(t.Value)-4:]
    if c.ColorMode {
        return "\033[33m" + masked + "\033[0m"
    }
    return masked
}
```

### fmt.Stringer / error

For non-struct types, `fmt.Stringer` and `error` implementations are used automatically.
For structs, field expansion takes priority (implement `PrettyPrinter` to override).

**Interface priority:**

1. `PrettyPrinterConfig` (config-aware)
2. `PrettyPrinter`
3. `fmt.Stringer` (non-struct only)
4. `error` (non-struct only)
5. Reflection-based formatting

## DefaultConfig

You can modify the global configuration:

```go
pf.DefaultConfig.ColorMode = false    // for logging
pf.DefaultConfig.UseJSONTags = true   // display with JSON names
pf.DefaultConfig.ShowTypes = true     // show type names
```

## License

MIT
