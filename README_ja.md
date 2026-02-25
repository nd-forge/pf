# pf

[English](README.md) | **日本語** | [中文](README_zh.md) | [한국어](README_ko.md) | [Español](README_es.md) | [Português](README_pt.md)

`fmt` を拡張した Go の pretty-print パッケージ。構造体・map・slice をインデント付きで見やすく出力します。

## インストール

```bash
go get github.com/nd-forge/pf
```

## クイックスタート

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

出力（ターミナルでは ANSI カラー付き）:
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

| 関数 | 説明 |
|---|---|
| `pf.Print(v)` | stdout に出力 |
| `pf.Sprint(v)` | string で返す |
| `pf.Fprint(w, v)` | io.Writer に出力 |

### Diff

| 関数 | 説明 |
|---|---|
| `pf.Diff(a, b)` | 差分を stdout に出力 |
| `pf.SprintDiff(a, b)` | 差分を string で返す |
| `pf.FprintDiff(w, a, b)` | 差分を io.Writer に出力 |

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
    Indent:      "    ",  // 4スペースインデント
    ShowTypes:   true,    // 型名を表示
    UseJSONTags: true,    // `json:"..."` タグ名を使用
    MaxDepth:    3,       // ネスト深度を制限
    ColorMode:   false,   // ANSIカラー無効（ログ向け）
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
// (Email は omitempty かつゼロ値 → 省略)
```

## インターフェース

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

Config を受け取ってフォーマットを変えられます。

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
構造体の場合はフィールド展開が優先されます（`PrettyPrinter` を実装すればそちらが優先）。

**インターフェース優先順位:**

1. `PrettyPrinterConfig`（Config 対応）
2. `PrettyPrinter`
3. `fmt.Stringer`（構造体以外）
4. `error`（構造体以外）
5. リフレクションベースのフォーマット

## DefaultConfig

グローバル設定を変更可能:

```go
pf.DefaultConfig.ColorMode = false    // ログ向け
pf.DefaultConfig.UseJSONTags = true   // JSON名で表示
pf.DefaultConfig.ShowTypes = true     // 型名表示
```

## License

MIT
