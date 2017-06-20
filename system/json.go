package system

import (
	"errors"
	"strconv"
	"strings"
)

// ToURL 将 interface{} 类型转成 URL
func ToURL(t interface{}) string {
	if t == nil {
		return ""
	}
	url := t.(string)
	if url[:4] != "http" {
		return "https:" + url
	}

	return url
}

// ToFloat 将 interface{} 类型转成 浮点类型
func ToFloat(t interface{}) float64 {
	var i float64

	if t == nil {
		return 0
	}
	t = indirect(t)
	switch val := t.(type) {
	case string:
		idx := strings.Index(val, "万")
		if idx >= 0 {
			if f, err := strconv.ParseFloat(val[:idx], 64); err == nil {
				i = f * 10000.0
			}
		} else {
			idx := strings.Index(val, "%")
			if idx >= 0 {
				if f, err := strconv.ParseFloat(val[:idx], 64); err == nil {
					i = f / 100.0
				}
			} else {
				i, _ = strconv.ParseFloat(val, 64)
			}
		}
	case float64:
		i = val
	case uint64:
		i = float64(val)
	case int64:
		i = float64(val)
	case int:
		i = float64(val)
	}

	return i
}

// ToInt 将 interface{} 类型转成整数
func ToInt(t interface{}) int {
	return int(ToFloat(t))
}

// ToString 将 interface{} 类型转成字符串
// func ToString(t interface{}) string {
// 	if t != nil {
// 		return t.(string)
// 	}

// 	return ""
// }

// GetSub 从 interface{} 取出子 interface{}
func GetSub(value interface{}, keys ...string) (interface{}, error) {
	var ok bool
	for _, k := range keys {
		values := value.(map[string]interface{})
		if value, ok = values[k]; ok == false {
			return nil, errors.New("No found")
		}
	}

	return value, nil
}

// PrintInterface 显示 Interface 接口
func PrintInterface(i interface{}) {
	s := StructToString(i)
	JSONPrint("", s)
}
