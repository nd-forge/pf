# pf

[English](README.md) | [日本語](README_ja.md) | [中文](README_zh.md) | **한국어** | [Español](README_es.md) | [Português](README_pt.md)

`fmt`를 확장한 Go pretty-print 패키지. 구조체, map, slice를 들여쓰기하여 보기 좋게 출력합니다.

## 설치

```bash
go get github.com/nd-forge/pf
```

## 빠른 시작

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

출력 (터미널에서 ANSI 컬러 적용):
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

| 함수 | 설명 |
|---|---|
| `pf.Print(v)` | stdout에 출력 |
| `pf.Sprint(v)` | string으로 반환 |
| `pf.Fprint(w, v)` | io.Writer에 출력 |

### Diff

| 함수 | 설명 |
|---|---|
| `pf.Diff(a, b)` | 차이를 stdout에 출력 |
| `pf.SprintDiff(a, b)` | 차이를 string으로 반환 |
| `pf.FprintDiff(w, a, b)` | 차이를 io.Writer에 출력 |

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
    Indent:      "    ",  // 4칸 들여쓰기
    ShowTypes:   true,    // 타입명 표시
    UseJSONTags: true,    // `json:"..."` 태그명 사용
    MaxDepth:    3,       // 중첩 깊이 제한
    ColorMode:   false,   // ANSI 컬러 비활성화 (로깅용)
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
// (Email은 omitempty이고 제로값 → 생략)
```

## 인터페이스

### PrettyPrinter

최우선. 타입에 대한 커스텀 pretty-print 출력을 정의합니다.

```go
type Token struct { Value string }

func (t Token) PrettyPrint() string {
    return fmt.Sprintf("Token(***%s)", t.Value[len(t.Value)-4:])
}

pf.Print(Token{Value: "secret1234"})
// Token(***1234)
```

### PrettyPrinterConfig

Config를 받아 포맷에 맞는 출력을 생성합니다.

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

구조체가 아닌 타입의 경우, `fmt.Stringer`와 `error` 구현이 자동으로 사용됩니다.
구조체의 경우 필드 전개가 우선됩니다 (`PrettyPrinter`를 구현하면 이를 재정의).

**인터페이스 우선순위:**

1. `PrettyPrinterConfig` (Config 인식)
2. `PrettyPrinter`
3. `fmt.Stringer` (구조체 외)
4. `error` (구조체 외)
5. 리플렉션 기반 포맷팅

## DefaultConfig

글로벌 설정을 변경할 수 있습니다:

```go
pf.DefaultConfig.ColorMode = false    // 로깅용
pf.DefaultConfig.UseJSONTags = true   // JSON 이름으로 표시
pf.DefaultConfig.ShowTypes = true     // 타입명 표시
```

## License

MIT
