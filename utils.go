package myGoUtils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

type Stringer interface {
	ToString() string
}

type number interface {
	integer | uinteger
}

type integer interface {
	int | int8 | int16 | int32 | int64
}

type uinteger interface {
	uint | uint8 | uint16 | uint32 | uint64
}

type double interface {
	float32 | float64
}

// TernaryOperator simulates the ternary operator just like (condition ? true : false)
func TernaryOperator(condition bool, ret1, ret2 any) any {
	if condition {
		return ret1
	}
	return ret2
}

// GetTagValue returns a value of a tag into a struct
func GetTagValue[Any any](fieldName, tagName string) (string, error) {
	var data Any
	dataType := reflect.TypeOf(data)

	if dataType.Kind() != reflect.Struct {
		return "", fmt.Errorf("'%s' must be a struct", dataType.String())
	}

	field, found := dataType.FieldByName(fieldName)
	if !found {
		return "", fmt.Errorf("field `%s` does not exist into `%s` struct", fieldName, dataType.String())
	}

	ret, found := field.Tag.Lookup(tagName)
	if !found {
		return "", fmt.Errorf("tag `%s` does not exist into `%s` field from `%s` struct", tagName, fieldName, dataType.String())
	}

	return ret, nil
}

// SplitAfterRegex splits by regex without losing the delimiters
func SplitAfterRegex(rgx *regexp.Regexp, str string) (ret []string) {
	var l []int
	i, all := 0, rgx.FindAllStringIndex(str, -1)
	if all[0][0] != 0 {
		l = append(l, 0)
	}
	for _, one := range all {
		l = append(l, one...)
		if add := str[l[i]:l[i+1]]; len(add) > 0 {
			ret = append(ret, add)
		}
		i++
	}
	if all[len(all)-1][1] != len(str) {
		l = append(l, len(str))
	}
	for ; i < len(l)-1; i++ {
		if add := str[l[i]:l[i+1]]; len(add) > 0 {
			ret = append(ret, add)
		}
	}
	return ret
}

// GetLocale get the system language
func GetLocale() string {

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("powershell", "Get-Culture | select -exp Name")

	case "darwin":
		cmd = exec.Command("osascript", "-e", "user locale of (get system info)")

	case "linux":
		output, ok := os.LookupEnv("LANG")
		if ok {
			return strings.ReplaceAll(strings.TrimSpace(output), "_", "-")
		}
	}

	output, err := cmd.Output()
	if err != nil || string(output) == "" {
		return "en-US"
	}

	return strings.ReplaceAll(strings.TrimSpace(string(output)), "_", "-")
}

// AbsoluteFlatMap converts arrays with recursive sub-arrays into simple array
func AbsoluteFlatMap(list []interface{}) []interface{} {
	var ret []interface{}
	for _, o := range list {
		one, ok := o.([]interface{})
		if ok {
			ret = append(ret, AbsoluteFlatMap(one)...)
		} else {
			ret = append(ret, o)
		}
	}
	return ret
}

// GetFilesInto creates an array composed by every file into a folder, you can filter by file type
func GetFilesInto(path string, extentionFilter string) ([]os.FileInfo, error) {
	var ret []os.FileInfo
	return ret, filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, extentionFilter) {
			ret = append(ret, info)
		}
		return nil
	})
}

// Int converts any valid numeral string into number
func Int[intType number](str string) (intType, error) {
	ret, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
	if err != nil {
		return 0, err
	}
	return intType(ret), nil
}

// Float converts any valid floating point string to a real floating point
func Float[floatType double](str string) (floatType, error) {
	ret, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
		return 0, nil
	}
	return floatType(ret), nil
}

// String converts anything to `string`
func String(v any) string {
	r := reflect.ValueOf(v)
	if r.Kind() == reflect.Struct {
		return "Object " + stringifyObj(r, 0, r.Kind())
	}
	return fmt.Sprintf("%v", v)
}

// Ptr get any obj pointer
func Ptr[Any any](some Any) *Any {
	return &some
}

// PtrVal returns the value held by a pointer. If the pointer is nil, a zero value is returned.
func PtrVal[Any any](some *Any) Any {
	if some == nil {
		var ret Any
		return ret
	}
	return *some
}

// MapKeys returns all keys of the map as a slice.
func MapKeys[K comparable, V any](m map[K]V) []K {
	k, _ := toSlice(m)
	return k
}

// MapValues returns all values of the map as a slice.
func MapValues[K comparable, V any](m map[K]V) []V {
	_, v := toSlice(m)
	return v
}

// ToPairs converts the map into a slice of key-value pairs.
func ToPairs[K comparable, V any](m map[K]V) []Pair[K, V] {
	i, pairs := 0, make([]Pair[K, V], len(m))
	for k, v := range m {
		pairs[i] = NewPair(k, v)
		i++
	}
	return pairs
}

// VecMap applies a function to each element in the slice and returns a new slice with the results
func VecMap[T, R any](v []T, f func(T) R) []R {
	count := len(v)
	if count == 0 {
		return nil
	}
	ret := make([]R, count)
	for i, value := range v {
		ret[i] = f(value)
	}
	return ret
}

// VecReduce reduces the slice to a single value using the provided function.
func VecReduce[T any](v []T, f func(T, T) T) (T, error) {
	count := len(v)
	if count == 0 {
		var ret T
		return ret, errors.New("empty array")
	}
	ret := v[0]
	for i := 1; i < count; i++ {
		ret = f(ret, v[i])
	}
	return ret, nil
}

// VecFilter separates elements in a slice based on a provided condition function.
// The first returned slice contains elements that satisfy the condition (approved),
// while the second slice contains elements that do not meet the condition (non-approved).
func VecFilter[T any](v []T, f func(T) bool) ([]T, []T) {
	count := len(v)
	if count == 0 {
		return nil, nil
	}
	i, ret, neg := 0, make([]T, 0, count), make([]T, 0, count)
	for _, value := range v {
		if f(value) {
			ret = append(ret, value)
			i++
			continue
		}
		neg = append(neg, value)
	}
	return ret, neg
}
