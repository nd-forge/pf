# pf

[English](README.md) | [日本語](README_ja.md) | [中文](README_zh.md) | [한국어](README_ko.md) | [Español](README_es.md) | **Português**

Um pacote Go de pretty-print que estende o `fmt`. Formata structs, maps e slices com indentação para uma saída legível.

## Instalação

```bash
go get github.com/nd-forge/pf
```

## Início rápido

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

Saída (com cores ANSI no terminal):
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

| Função | Descrição |
|---|---|
| `pf.Print(v)` | Imprime no stdout |
| `pf.Sprint(v)` | Retorna como string |
| `pf.Fprint(w, v)` | Escreve no io.Writer |

### Diff

| Função | Descrição |
|---|---|
| `pf.Diff(a, b)` | Imprime as diferenças no stdout |
| `pf.SprintDiff(a, b)` | Retorna as diferenças como string |
| `pf.FprintDiff(w, a, b)` | Escreve as diferenças no io.Writer |

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
    Indent:      "    ",  // indentação de 4 espaços
    ShowTypes:   true,    // mostrar nomes de tipo
    UseJSONTags: true,    // usar nomes das tags `json:"..."`
    MaxDepth:    3,       // limitar profundidade de aninhamento
    ColorMode:   false,   // sem cores ANSI (para logs)
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
// (Email é omitempty e valor zero → omitido)
```

## Interfaces

### PrettyPrinter

Maior prioridade. Defina uma saída pretty-print personalizada para seus tipos.

```go
type Token struct { Value string }

func (t Token) PrettyPrint() string {
    return fmt.Sprintf("Token(***%s)", t.Value[len(t.Value)-4:])
}

pf.Print(Token{Value: "secret1234"})
// Token(***1234)
```

### PrettyPrinterConfig

Recebe o Config para produzir uma saída adaptada ao formato.

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

Para tipos que não são structs, as implementações de `fmt.Stringer` e `error` são usadas automaticamente.
Para structs, a expansão de campos tem prioridade (implemente `PrettyPrinter` para substituir).

**Prioridade de interfaces:**

1. `PrettyPrinterConfig` (com reconhecimento de config)
2. `PrettyPrinter`
3. `fmt.Stringer` (apenas não-struct)
4. `error` (apenas não-struct)
5. Formatação baseada em reflexão

## DefaultConfig

Você pode modificar a configuração global:

```go
pf.DefaultConfig.ColorMode = false    // para logs
pf.DefaultConfig.UseJSONTags = true   // exibir com nomes JSON
pf.DefaultConfig.ShowTypes = true     // mostrar nomes de tipo
```

## License

MIT
