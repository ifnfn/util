package system

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/cast"
)

func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

func arrayToString(a []interface{}) string {
	out := make([]string, len(a))
	for i, elm := range a {
		out[i] = ToString(elm)
	}
	return fmt.Sprintf("%s", strings.Join(out, ", "))
}

func intToString(a []uint64) string {
	out := make([]string, len(a))
	for i, elm := range a {
		out[i] = strconv.FormatUint(elm, 10)
	}
	return fmt.Sprintf("%s", strings.Join(out, ", "))
}

func strArrayToString(a []string) string {
	out := make([]string, len(a))
	for i, elm := range a {
		out[i] = elm
	}
	return fmt.Sprintf("%s", strings.Join(out, ", "))
}

func strMapToString(m map[string]interface{}) string {
	out := make([]string, len(m))
	i := 0
	for key, value := range m {
		out[i] = fmt.Sprintf("%s: %s", key, ToString(value))
		i++
	}
	return fmt.Sprintf("{%s}", strings.Join(out, ", "))
}

func mapToString(m map[interface{}]interface{}) string {
	out := make([]string, len(m))
	i := 0
	for key, value := range m {
		out[i] = fmt.Sprintf("%s: %s", ToString(key), ToString(value))
		i++
	}
	return fmt.Sprintf("{%s}", strings.Join(out, ", "))
}

// ToString 将多种类型转成字符串
func ToString(i interface{}) string {
	switch val := i.(type) {
	case error:
		return val.Error()
	}
	i = indirect(i)
	switch val := i.(type) {
	case uint64:
		return strconv.FormatUint(val, 10)
	case float64:
		return strconv.FormatFloat(val, 'f', 0, 64)
	case []interface{}:
		return arrayToString(val)
	case []string:
		return strArrayToString(val)
	case []uint64:
		return intToString(val)
	case map[string]interface{}:
		return strMapToString(val)
	case map[interface{}]interface{}:
		return mapToString(val)
	default:
		newval, err := cast.ToStringE(val)
		if err != nil {
			return fmt.Sprint(newval)
		}
		return newval
	}
}

// StringArrayToString 将字符串数组转成字符串
func StringArrayToString(a []string, quote string, separator string) string {
	out := make([]string, len(a))
	for i, elm := range a {
		out[i] = quote + elm + quote
	}
	return fmt.Sprintf("%s", strings.Join(out, separator))
}
