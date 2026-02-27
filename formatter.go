package pf

import (
	"fmt"
	"reflect"
	"strings"
)

func (f *formatter) format(v reflect.Value, depth int) {
	if f.config.MaxDepth > 0 && depth > f.config.MaxDepth {
		f.sb.WriteString("...")
		return
	}

	// Handle invalid (nil interface)
	if !v.IsValid() {
		f.colored(cNil, "nil")
		return
	}

	// Check interfaces BEFORE dereferencing pointers,
	// so pointer receivers work too.
	if f.tryInterfaces(v) {
		return
	}

	// Dereference pointers
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			f.colored(cNil, "nil")
			return
		}
		v = v.Elem()
		// Check interfaces again on the dereferenced value
		if f.tryInterfaces(v) {
			return
		}
	}

	f.formatByKind(v, depth)
}

func (f *formatter) formatByKind(v reflect.Value, depth int) {
	switch v.Kind() {
	case reflect.Struct:
		f.formatStruct(v, depth)
	case reflect.Map:
		f.formatMap(v, depth)
	case reflect.Slice, reflect.Array:
		f.formatSlice(v, depth)
	case reflect.String:
		f.colored(cString, fmt.Sprintf("%q", v.String()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		f.colored(cNumber, fmt.Sprintf("%d", v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		f.colored(cNumber, fmt.Sprintf("%d", v.Uint()))
	case reflect.Float32, reflect.Float64:
		f.colored(cNumber, formatFloat(v.Float()))
	case reflect.Bool:
		f.colored(cBool, fmt.Sprintf("%t", v.Bool()))
	case reflect.Interface:
		if v.IsNil() {
			f.colored(cNil, "nil")
		} else {
			f.format(v.Elem(), depth)
		}
	case reflect.Chan:
		f.colored(cType, fmt.Sprintf("(chan %s)", v.Type().Elem()))
	case reflect.Func:
		f.colored(cType, fmt.Sprintf("(func %s)", v.Type()))
	default:
		f.sb.WriteString(fmt.Sprintf("%v", v.Interface()))
	}
}

// tryInterfaces checks if the value implements PrettyPrinterConfig,
// PrettyPrinter, fmt.Stringer, or error, in that order.
// Returns true if an interface was used to format the value.
func (f *formatter) tryInterfaces(v reflect.Value) bool {
	if !v.IsValid() || !v.CanInterface() {
		return false
	}

	iface := v.Interface()

	// 1. PrettyPrinterConfig (highest priority — format-aware)
	if pp, ok := iface.(PrettyPrinterConfig); ok {
		f.sb.WriteString(pp.PrettyPrintConfig(f.config))
		return true
	}

	// 2. PrettyPrinter
	if pp, ok := iface.(PrettyPrinter); ok {
		f.sb.WriteString(pp.PrettyPrint())
		return true
	}

	// Also check pointer to value for pointer receiver methods
	if v.Kind() != reflect.Ptr && v.CanAddr() {
		ptr := v.Addr().Interface()
		if pp, ok := ptr.(PrettyPrinterConfig); ok {
			f.sb.WriteString(pp.PrettyPrintConfig(f.config))
			return true
		}
		if pp, ok := ptr.(PrettyPrinter); ok {
			f.sb.WriteString(pp.PrettyPrint())
			return true
		}
	}

	// 3. fmt.Stringer — only for non-struct types to avoid
	//    losing struct detail (many structs implement Stringer
	//    but you still want to see inside them by default)
	if v.Kind() != reflect.Struct {
		if s, ok := iface.(fmt.Stringer); ok {
			f.colored(cString, fmt.Sprintf("%q", s.String()))
			return true
		}
	}

	// 4. error interface
	if v.Kind() != reflect.Struct {
		if e, ok := iface.(error); ok {
			f.colored(cNil, fmt.Sprintf("error(%q)", e.Error()))
			return true
		}
	}

	return false
}

func (f *formatter) formatStruct(v reflect.Value, depth int) {
	t := v.Type()
	indent := strings.Repeat(f.config.Indent, depth+1)
	closingIndent := strings.Repeat(f.config.Indent, depth)

	if f.config.ShowTypes {
		f.colored(cType, t.Name()+" ")
	}

	// Collect visible fields
	type fieldEntry struct {
		displayName string
		value       reflect.Value
	}
	var fields []fieldEntry

	for i := 0; i < v.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}

		name := resolveFieldName(sf, v.Field(i), f.config.UseJSONTags)
		if name == "" {
			continue
		}

		fields = append(fields, fieldEntry{
			displayName: name,
			value:       v.Field(i),
		})
	}

	if len(fields) == 0 {
		f.colored(cBrace, "{}")
		return
	}

	f.colored(cBrace, "{\n")

	for i, fe := range fields {
		f.sb.WriteString(indent)
		f.colored(cKey, fe.displayName)
		f.sb.WriteString(": ")
		f.format(fe.value, depth+1)
		if i < len(fields)-1 {
			f.sb.WriteString(",")
		}
		f.sb.WriteString("\n")
	}

	f.sb.WriteString(closingIndent)
	f.colored(cBrace, "}")
}

// resolveFieldName returns the display name for a struct field.
// Returns "" if the field should be skipped.
func resolveFieldName(sf reflect.StructField, fieldValue reflect.Value, useJSONTags bool) string {
	name := sf.Name
	if !useJSONTags {
		return name
	}

	tag := sf.Tag.Get("json")
	if tag == "" {
		return name
	}

	parts := strings.Split(tag, ",")
	if parts[0] == "-" {
		return ""
	}
	if parts[0] != "" {
		name = parts[0]
	}

	for _, opt := range parts[1:] {
		if opt == "omitempty" && fieldValue.IsZero() {
			return ""
		}
	}
	return name
}

func (f *formatter) formatMap(v reflect.Value, depth int) {
	if v.IsNil() {
		f.colored(cNil, "nil")
		return
	}

	indent := strings.Repeat(f.config.Indent, depth+1)
	closingIndent := strings.Repeat(f.config.Indent, depth)

	if f.config.ShowTypes {
		f.colored(cType, fmt.Sprintf("map[%s]%s ", v.Type().Key(), v.Type().Elem()))
	}

	keys := v.MapKeys()
	if len(keys) == 0 {
		f.colored(cBrace, "{}")
		return
	}

	// Sort keys for deterministic output
	sortedKeys := sortMapKeys(keys)

	f.colored(cBrace, "{\n")
	for i, key := range sortedKeys {
		f.sb.WriteString(indent)
		f.format(key, depth+1)
		f.sb.WriteString(": ")
		f.format(v.MapIndex(key), depth+1)
		if i < len(sortedKeys)-1 {
			f.sb.WriteString(",")
		}
		f.sb.WriteString("\n")
	}
	f.sb.WriteString(closingIndent)
	f.colored(cBrace, "}")
}

func (f *formatter) formatSlice(v reflect.Value, depth int) {
	if v.Kind() == reflect.Slice && v.IsNil() {
		f.colored(cNil, "nil")
		return
	}

	if v.Len() == 0 {
		f.colored(cBrace, "[]")
		return
	}

	// Compact for short simple slices
	if v.Len() <= 5 && isSimpleKind(v.Type().Elem().Kind()) {
		f.colored(cBrace, "[")
		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				f.sb.WriteString(", ")
			}
			f.format(v.Index(i), depth)
		}
		f.colored(cBrace, "]")
		return
	}

	indent := strings.Repeat(f.config.Indent, depth+1)
	closingIndent := strings.Repeat(f.config.Indent, depth)

	f.colored(cBrace, "[\n")
	for i := 0; i < v.Len(); i++ {
		f.sb.WriteString(indent)
		f.format(v.Index(i), depth+1)
		if i < v.Len()-1 {
			f.sb.WriteString(",")
		}
		f.sb.WriteString("\n")
	}
	f.sb.WriteString(closingIndent)
	f.colored(cBrace, "]")
}

// --- Helpers ---

func isSimpleKind(k reflect.Kind) bool {
	switch k {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return true
	}
	return false
}

func formatFloat(f float64) string {
	// Use comma-friendly format: show decimals only if needed
	if f == float64(int64(f)) {
		return fmt.Sprintf("%.1f", f)
	}
	return fmt.Sprintf("%g", f)
}

// sortMapKeys sorts reflect.Value keys by their string representation.
func sortMapKeys(keys []reflect.Value) []reflect.Value {
	sorted := make([]reflect.Value, len(keys))
	copy(sorted, keys)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if fmt.Sprint(sorted[i].Interface()) > fmt.Sprint(sorted[j].Interface()) {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}
