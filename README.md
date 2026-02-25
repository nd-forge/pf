# pf

`fmt` を拡張した Go の pretty-print パッケージ。構造体・map・slice をインデント付きで見やすく出力します。

## Install

```bash
go get github.com/taiki-nd/pf
```

## Quick Start

```go
import "github.com/taiki-nd/pf"

type User struct {
    Name    string `json:"name"`
    Age     int    `json:"age"`
    Active  bool
    Address Address
}

user := User{Name: "野田", Age: 30, Active: true, Address: Address{City: "横浜"}}

// Pretty print
pf.Print(user)
```

Output (with ANSI colors in terminal):
```
{
  Name: "野田",
  Age: 30,
  Active: true,
  Address: {
    City: "横浜"
  }
}
```

## API

### Pretty Print

| Function | Description |
|---|---|
| `pf.Print(v)` | stdout に出力 |
| `pf.Sprint(v)` | string で返す |
| `pf.Fprint(w, v)` | io.Writer に出力 |

### Diff

| Function | Description |
|---|---|
| `pf.Diff(a, b)` | 差分を stdout に出力 |
| `pf.SprintDiff(a, b)` | 差分を string で返す |
| `pf.FprintDiff(w, a, b)` | 差分を io.Writer に出力 |

```go
old := User{Name: "Noda", Age: 30, Active: true}
new := User{Name: "Noda", Age: 31, Active: false}

pf.Diff(old, new)
// {
//   Name: "Noda"
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
//   user_name: "野田",
// }
// (Email is omitempty + zero value → omitted)
```

## Interfaces

### PrettyPrinter

最優先。独自の pretty-print 出力を定義できます。

```go
type Token struct { Value string }

func (t Token) PrettyPrint() string {
    return fmt.Sprintf("Token(***%s)", t.Value[len(t.Value)-4:])
}

pf.Print(Token{Value: "secret1234"})
// Token(***1234)
```

### PrettyPrinterConfig

Config を受け取って出力を変えられます。

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

構造体以外の型で `fmt.Stringer` や `error` を実装していれば自動的に使われます。
構造体の場合はフィールド展開が優先されます（PrettyPrinter を実装すればそちらが優先）。

**Interface 優先順位:**

1. `PrettyPrinterConfig` (config-aware)
2. `PrettyPrinter`
3. `fmt.Stringer` (non-struct only)
4. `error` (non-struct only)
5. Reflection-based formatting

## DefaultConfig

グローバル設定を変更可能：

```go
pf.DefaultConfig.ColorMode = false    // ログ向け
pf.DefaultConfig.UseJSONTags = true   // JSON名で表示
pf.DefaultConfig.ShowTypes = true     // 型名表示
```

## License

MIT
