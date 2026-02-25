# pf

[English](README.md) | [日本語](README_ja.md) | **中文** | [한국어](README_ko.md) | [Español](README_es.md) | [Português](README_pt.md)

扩展 `fmt` 的 Go pretty-print 包。以缩进格式美观地输出结构体、map 和 slice。

## 安装

```bash
go get github.com/nd-forge/pf
```

## 快速开始

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

输出（终端中带 ANSI 颜色）:
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

| 函数 | 说明 |
|---|---|
| `pf.Print(v)` | 输出到 stdout |
| `pf.Sprint(v)` | 返回 string |
| `pf.Fprint(w, v)` | 输出到 io.Writer |

### Diff

| 函数 | 说明 |
|---|---|
| `pf.Diff(a, b)` | 将差异输出到 stdout |
| `pf.SprintDiff(a, b)` | 返回差异 string |
| `pf.FprintDiff(w, a, b)` | 将差异输出到 io.Writer |

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
    Indent:      "    ",  // 4空格缩进
    ShowTypes:   true,    // 显示类型名
    UseJSONTags: true,    // 使用 `json:"..."` 标签名
    MaxDepth:    3,       // 限制嵌套深度
    ColorMode:   false,   // 禁用 ANSI 颜色（适用于日志）
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
// (Email 为 omitempty 且零值 → 省略)
```

## 接口

### PrettyPrinter

最高优先级。为类型定义自定义的 pretty-print 输出。

```go
type Token struct { Value string }

func (t Token) PrettyPrint() string {
    return fmt.Sprintf("Token(***%s)", t.Value[len(t.Value)-4:])
}

pf.Print(Token{Value: "secret1234"})
// Token(***1234)
```

### PrettyPrinterConfig

接收 Config 以生成格式感知的输出。

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

对于非结构体类型，会自动使用 `fmt.Stringer` 和 `error` 的实现。
对于结构体，字段展开优先（实现 `PrettyPrinter` 可覆盖此行为）。

**接口优先级:**

1. `PrettyPrinterConfig`（感知配置）
2. `PrettyPrinter`
3. `fmt.Stringer`（仅非结构体）
4. `error`（仅非结构体）
5. 基于反射的格式化

## DefaultConfig

可修改全局配置:

```go
pf.DefaultConfig.ColorMode = false    // 适用于日志
pf.DefaultConfig.UseJSONTags = true   // 使用 JSON 名称显示
pf.DefaultConfig.ShowTypes = true     // 显示类型名
```

## License

MIT
