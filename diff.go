package pf

import (
	"fmt"
	"reflect"
	"strings"
)

type differ struct {
	config Config
	sb     strings.Builder
}

// diff compares two values and returns a formatted diff string.
// For structs, it shows changed fields with -/+ markers.
// For non-structs, it shows a simple before/after.
func (d *differ) diff(a, b interface{}) string {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	// Dereference pointers
	for va.Kind() == reflect.Ptr && !va.IsNil() {
		va = va.Elem()
	}
	for vb.Kind() == reflect.Ptr && !vb.IsNil() {
		vb = vb.Elem()
	}

	if va.Type() != vb.Type() {
		d.writeLine("", fmt.Sprintf("type mismatch: %s vs %s", va.Type(), vb.Type()))
		return d.sb.String()
	}

	switch va.Kind() {
	case reflect.Struct:
		d.diffStruct(va, vb, 0)
	case reflect.Map:
		d.diffMap(va, vb, 0)
	case reflect.Slice, reflect.Array:
		d.diffSlice(va, vb, 0)
	default:
		d.diffScalar(va, vb, 0)
	}

	return d.sb.String()
}

func (d *differ) diffStruct(a, b reflect.Value, depth int) {
	t := a.Type()
	indent := strings.Repeat(d.config.Indent, depth+1)
	closingIndent := strings.Repeat(d.config.Indent, depth)
	cm := d.config.ColorMode

	typeName := t.Name()
	if d.config.ShowTypes && typeName != "" {
		d.sb.WriteString(coloredStr(cType, typeName+" ", cm))
	}
	d.sb.WriteString(coloredStr(cBrace, "{\n", cm))

	for i := 0; i < a.NumField(); i++ {
		sf := t.Field(i)
		if !sf.IsExported() {
			continue
		}

		name := d.fieldName(sf)
		if name == "" {
			continue
		}

		fa := a.Field(i)
		fb := b.Field(i)

		aStr := d.sprintValue(fa)
		bStr := d.sprintValue(fb)

		if aStr == bStr {
			// Unchanged
			d.sb.WriteString(indent)
			d.sb.WriteString(coloredStr(cKey, name, cm))
			d.sb.WriteString(": ")
			d.sb.WriteString(aStr)
			d.sb.WriteString("\n")
		} else {
			// Changed
			d.sb.WriteString(indent)
			d.sb.WriteString(coloredStr(cDiffDel, "- "+name+": "+aStr, cm))
			d.sb.WriteString("\n")
			d.sb.WriteString(indent)
			d.sb.WriteString(coloredStr(cDiffAdd, "+ "+name+": "+bStr, cm))
			d.sb.WriteString("\n")
		}
	}

	d.sb.WriteString(closingIndent)
	d.sb.WriteString(coloredStr(cBrace, "}", cm))
}

func (d *differ) diffMap(a, b reflect.Value, depth int) {
	indent := strings.Repeat(d.config.Indent, depth+1)
	closingIndent := strings.Repeat(d.config.Indent, depth)
	cm := d.config.ColorMode

	d.sb.WriteString(coloredStr(cBrace, "{\n", cm))

	// Collect all keys from both maps
	allKeys := make(map[string]reflect.Value)
	for _, k := range a.MapKeys() {
		allKeys[fmt.Sprint(k.Interface())] = k
	}
	for _, k := range b.MapKeys() {
		allKeys[fmt.Sprint(k.Interface())] = k
	}

	for _, key := range sortMapKeys(collectValues(allKeys)) {
		keyStr := fmt.Sprint(key.Interface())
		aVal := a.MapIndex(key)
		bVal := b.MapIndex(key)

		aExists := aVal.IsValid()
		bExists := bVal.IsValid()

		switch {
		case aExists && !bExists:
			d.sb.WriteString(indent)
			d.sb.WriteString(coloredStr(cDiffDel, "- "+keyStr+": "+d.sprintValue(aVal), cm))
			d.sb.WriteString("\n")
		case !aExists && bExists:
			d.sb.WriteString(indent)
			d.sb.WriteString(coloredStr(cDiffAdd, "+ "+keyStr+": "+d.sprintValue(bVal), cm))
			d.sb.WriteString("\n")
		default:
			aStr := d.sprintValue(aVal)
			bStr := d.sprintValue(bVal)
			if aStr == bStr {
				d.sb.WriteString(indent)
				d.sb.WriteString(keyStr + ": " + aStr + "\n")
			} else {
				d.sb.WriteString(indent)
				d.sb.WriteString(coloredStr(cDiffDel, "- "+keyStr+": "+aStr, cm))
				d.sb.WriteString("\n")
				d.sb.WriteString(indent)
				d.sb.WriteString(coloredStr(cDiffAdd, "+ "+keyStr+": "+bStr, cm))
				d.sb.WriteString("\n")
			}
		}
	}

	d.sb.WriteString(closingIndent)
	d.sb.WriteString(coloredStr(cBrace, "}", cm))
}

func (d *differ) diffSlice(a, b reflect.Value, depth int) {
	indent := strings.Repeat(d.config.Indent, depth+1)
	closingIndent := strings.Repeat(d.config.Indent, depth)
	cm := d.config.ColorMode

	d.sb.WriteString(coloredStr(cBrace, "[\n", cm))

	maxLen := a.Len()
	if b.Len() > maxLen {
		maxLen = b.Len()
	}

	for i := 0; i < maxLen; i++ {
		aExists := i < a.Len()
		bExists := i < b.Len()

		switch {
		case aExists && !bExists:
			d.sb.WriteString(indent)
			d.sb.WriteString(coloredStr(cDiffDel, fmt.Sprintf("- [%d]: %s", i, d.sprintValue(a.Index(i))), cm))
			d.sb.WriteString("\n")
		case !aExists && bExists:
			d.sb.WriteString(indent)
			d.sb.WriteString(coloredStr(cDiffAdd, fmt.Sprintf("+ [%d]: %s", i, d.sprintValue(b.Index(i))), cm))
			d.sb.WriteString("\n")
		default:
			aStr := d.sprintValue(a.Index(i))
			bStr := d.sprintValue(b.Index(i))
			if aStr == bStr {
				d.sb.WriteString(indent)
				d.sb.WriteString(fmt.Sprintf("[%d]: %s\n", i, aStr))
			} else {
				d.sb.WriteString(indent)
				d.sb.WriteString(coloredStr(cDiffDel, fmt.Sprintf("- [%d]: %s", i, aStr), cm))
				d.sb.WriteString("\n")
				d.sb.WriteString(indent)
				d.sb.WriteString(coloredStr(cDiffAdd, fmt.Sprintf("+ [%d]: %s", i, bStr), cm))
				d.sb.WriteString("\n")
			}
		}
	}

	d.sb.WriteString(closingIndent)
	d.sb.WriteString(coloredStr(cBrace, "]", cm))
}

func (d *differ) diffScalar(a, b reflect.Value, depth int) {
	cm := d.config.ColorMode
	aStr := d.sprintValue(a)
	bStr := d.sprintValue(b)
	if aStr == bStr {
		d.sb.WriteString(aStr)
	} else {
		d.sb.WriteString(coloredStr(cDiffDel, "- "+aStr, cm))
		d.sb.WriteString("\n")
		d.sb.WriteString(coloredStr(cDiffAdd, "+ "+bStr, cm))
	}
}

func (d *differ) sprintValue(v reflect.Value) string {
	noColor := Config{
		Indent:      d.config.Indent,
		ShowTypes:   d.config.ShowTypes,
		UseJSONTags: d.config.UseJSONTags,
		MaxDepth:    d.config.MaxDepth,
		ColorMode:   false, // no color for comparison
	}
	return noColor.Sprint(v.Interface())
}

func (d *differ) fieldName(sf reflect.StructField) string {
	if d.config.UseJSONTags {
		if tag := sf.Tag.Get("json"); tag != "" {
			parts := strings.Split(tag, ",")
			if parts[0] == "-" {
				return ""
			}
			if parts[0] != "" {
				return parts[0]
			}
		}
	}
	return sf.Name
}

func (d *differ) writeLine(prefix, text string) {
	if prefix != "" {
		d.sb.WriteString(prefix + " ")
	}
	d.sb.WriteString(text)
	d.sb.WriteString("\n")
}

func collectValues(m map[string]reflect.Value) []reflect.Value {
	vals := make([]reflect.Value, 0, len(m))
	for _, v := range m {
		vals = append(vals, v)
	}
	return vals
}
