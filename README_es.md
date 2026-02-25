# pf

[English](README.md) | [日本語](README_ja.md) | [中文](README_zh.md) | [한국어](README_ko.md) | **Español** | [Português](README_pt.md)

Un paquete Go de pretty-print que extiende `fmt`. Formatea structs, maps y slices con indentación para una salida legible.

## Instalación

```bash
go get github.com/nd-forge/pf
```

## Inicio rápido

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

Salida (con colores ANSI en la terminal):
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

| Función | Descripción |
|---|---|
| `pf.Print(v)` | Imprime en stdout |
| `pf.Sprint(v)` | Devuelve como string |
| `pf.Fprint(w, v)` | Escribe en io.Writer |

### Diff

| Función | Descripción |
|---|---|
| `pf.Diff(a, b)` | Imprime las diferencias en stdout |
| `pf.SprintDiff(a, b)` | Devuelve las diferencias como string |
| `pf.FprintDiff(w, a, b)` | Escribe las diferencias en io.Writer |

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
    Indent:      "    ",  // indentación de 4 espacios
    ShowTypes:   true,    // mostrar nombres de tipo
    UseJSONTags: true,    // usar nombres de etiquetas `json:"..."`
    MaxDepth:    3,       // limitar profundidad de anidamiento
    ColorMode:   false,   // sin colores ANSI (para logs)
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
// (Email es omitempty y valor cero → omitido)
```

## Interfaces

### PrettyPrinter

Máxima prioridad. Define una salida pretty-print personalizada para tus tipos.

```go
type Token struct { Value string }

func (t Token) PrettyPrint() string {
    return fmt.Sprintf("Token(***%s)", t.Value[len(t.Value)-4:])
}

pf.Print(Token{Value: "secret1234"})
// Token(***1234)
```

### PrettyPrinterConfig

Recibe el Config para producir una salida adaptada al formato.

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

Para tipos que no son structs, las implementaciones de `fmt.Stringer` y `error` se usan automáticamente.
Para structs, la expansión de campos tiene prioridad (implementa `PrettyPrinter` para anular esto).

**Prioridad de interfaces:**

1. `PrettyPrinterConfig` (con reconocimiento de config)
2. `PrettyPrinter`
3. `fmt.Stringer` (solo no-struct)
4. `error` (solo no-struct)
5. Formateo basado en reflexión

## DefaultConfig

Puedes modificar la configuración global:

```go
pf.DefaultConfig.ColorMode = false    // para logs
pf.DefaultConfig.UseJSONTags = true   // mostrar con nombres JSON
pf.DefaultConfig.ShowTypes = true     // mostrar nombres de tipo
```

## License

MIT
