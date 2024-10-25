package myGoUtils

import (
	"fmt"
	"reflect"
	"strings"
)

func toSlice[K comparable, V any](m map[K]V) ([]K, []V) {
	k := make([]K, len(m))
	i, v := 0, make([]V, len(m))
	for key, value := range m {
		k[i], v[i] = key, value
		i++
	}
	return k, v
}

func stringifyObj(v reflect.Value, indentLevel int, root reflect.Kind) string {
	const indentChar = "    "
	indent := strings.Repeat(indentChar, indentLevel)
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return "nil"
		}
		return "&" + stringifyObj(v.Elem(), indentLevel, root)
	case reflect.Struct:
		var result strings.Builder
		result.WriteString("{\n")
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			value := v.Field(i)
			structStr := stringifyObj(value, indentLevel+1, root)
			str := fmt.Sprintf("%s.%s %s\n", indent+indentChar, field.Name, structStr)
			result.WriteString(str)
		}
		result.WriteString(indent + "}")
		return result.String()
	case reflect.Slice:
		if v.Len() == 0 {
			return "[]"
		}
		var result strings.Builder
		result.WriteString("[\n")
		for i := 0; i < v.Len(); i++ {
			sliceStr := stringifyObj(v.Index(i), indentLevel+1, root)
			str := fmt.Sprintf("%s%s\n", indent+indentChar, sliceStr)
			result.WriteString(str)
		}
		result.WriteString(indent + "]")
		return result.String()
	case reflect.String:
		if root == reflect.String {
			return fmt.Sprintf("%s", v.String())
		}
		return fmt.Sprintf("\"%s\"", v.String())
	case reflect.Int:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Interface, reflect.Invalid:
		return "nil"
	default:
		return fmt.Sprintf("%v", v)
	}
}
