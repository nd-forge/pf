package pf

// PrettyPrinter can be implemented by any type to control its own
// pretty-print output. When a value implements PrettyPrinter,
// pf uses PrettyPrint() instead of reflection-based formatting.
//
//	type User struct { Name string; Age int }
//
//	func (u User) PrettyPrint() string {
//	    return fmt.Sprintf("User<%s, age=%d>", u.Name, u.Age)
//	}
type PrettyPrinter interface {
	PrettyPrint() string
}

// PrettyPrinterConfig is like PrettyPrinter but receives the current
// Config, allowing format-aware custom output.
//
//	func (u User) PrettyPrintConfig(c pf.Config) string {
//	    if c.ColorMode {
//	        return "\033[36mUser<" + u.Name + ">\033[0m"
//	    }
//	    return "User<" + u.Name + ">"
//	}
type PrettyPrinterConfig interface {
	PrettyPrintConfig(c Config) string
}
