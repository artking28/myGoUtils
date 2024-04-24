package myGoUtils

import (
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

type Atom interface {
	number | double | any | struct{} | string
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
func TernaryOperator[Any Atom, Any2 Atom](condition bool, ret1 Any, ret2 Any2) any {
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

// SplitAfterRegex splits by regex withot losing the delimiters
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

// AbsoluteFlatMap converts arrays with sub-arrays into a linear array line `[[[[][]][[[]]][]][]} => []`
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
func String[Any any](some Any) string {
	return fmt.Sprintf("%v", some)
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
