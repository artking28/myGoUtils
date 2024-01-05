package myGoUtils
 
import ( 
   "errors"
   "fmt"
   "path/filepath"
   "os"
   "reflect"
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

// TernaryOperator simulate the ternary operator just like (condition ? true : false)
func TernaryOpeator[Any Atom, Any2 Atom](condition bool, ret1 Any, ret2 Any2) any {
	if condition {
		return ret1
	}
   return ret2
}

// GetTagValue returns a value of a tag into a struct
func GetTagValue[Any any](fieldName, tagName string) (string, error) {
    // Obtém o tipo da struct
    var data Any
    dataType := reflect.TypeOf(data)

    // Verifica se o tipo é uma struct
    if dataType.Kind() != reflect.Struct {
        return "", fmt.Errorf("'%s' must be a struct", dataType.String())
    }

    // Obtém o valor do campo
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

// GetFilesInto create an array composed by every file into a folder, you can filter by file type
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

// Int convert any valid numeral string into number
func Int[intType number](str string) (intType, error) {
   ret, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64)
	if err != nil{
      return 0, err
	}
	return intType(ret), nil
}

// Float convert any valid floating point string to a real floating point
func Float[floatType double](str string) (floatType, error) {
   ret, err := strconv.ParseFloat(strings.TrimSpace(str), 64)
	if err != nil {
      return 0, nil
	}
	return floatType(ret), nil
}

// String convert anything to `string`
func String[Any any](some Any) string {
	return fmt.Sprintf("%v", some)
}

// Ptr get any obj pointer
func Ptr[Any any](some Any) *Any {
	return &some
}

/* PointerVal returns the value held by a pointer. If the pointer
 is nil the zerovalue is returned. */
func PtrVal[Any any](some *Any) Any {
	if some == nil {
		var ret Any
		return ret
	}
	return *some
}
