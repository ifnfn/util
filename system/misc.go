package system

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetTime 将当前时间转成字符串
func GetTime() string {
	t := time.Now() //获取当前时间的结构体
	return fmt.Sprintf(t.Format("2006-01-02 15:04:05"))
}

// GetTimeInt 当前时间 Unix 值
func GetTimeInt() int64 {
	return time.Now().Unix()
}

// ConvTimeToInt string time to int64
func ConvTimeToInt(strTime string) int64 {
	tm, _ := time.ParseInLocation("2006-01-02 15:04:05", strTime, time.Local)
	return tm.Unix()
}

// ConvTimeToString 将int64 时间转成字符串
func ConvTimeToString(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02 15:04:05")
}

// GetAliyunFormatTime 阿里云格式时间？
func GetAliyunFormatTime() string {
	ow := time.Now()
	year, mon, day := ow.UTC().Date()
	hour, min, sec := ow.UTC().Clock()
	strTime := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ", year, mon, day, hour, min, sec)
	return strTime
}

// GetMD5 计算字符串的MD5字符串
func GetMD5(text []byte) string {
	sum := md5.Sum(text)

	return hex.EncodeToString(sum[:])
}

// Md5File 计算文件的 MD5值字符串
func Md5File(strPath string) string {
	file, err := os.Open(strPath)
	if err != nil {
		return ""
	}
	defer file.Close()

	h := md5.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(h.Sum(nil))
}

// GetCurrentDirectory 获得程序当前目录
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Fatal(err)
	}

	return strings.Replace(dir, "\\", "/", -1)
}

// StructToString 将结构体转成字符串
func StructToString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case int:
		return fmt.Sprintf("%d", t)
	default:
		b, err := json.Marshal(v)
		if err != nil {
			logger.Warn(err)
		}
		str := string(b)

		return str
	}
}

//截取字符串末几位
func SubRightstr(str string, length int) string {
	n := len(str)
	if n < length {
		return ""
	}
	return Substr(str, n-length, length)
}

// Substr 截取字符串 start 起点下标 length 需要截取的长度
func Substr(str string, start int, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}

	return string(rs[start:end])
}

// Substr2 截取字符串 start 起点下标 end 终点下标(不包括)
func Substr2(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

// JSONPrint JSON 格式化打印
func JSONPrint(name, val string) {
	if true {
		fmt.Println(name)
		var out bytes.Buffer
		err := json.Indent(&out, []byte(val), "", "  ")

		if err != nil {
			println(err.Error())
		}

		out.WriteTo(os.Stdout)
		println()
	}
}

// JSON2Map JSON 转 MAP（仅支持一纬转换）
func JSON2Map(jsonData string) (s map[string]interface{}, err error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// MapPrint Map =>String=>Print  格式化打印
func MapPrint(m interface{}) {
	//	b, _ := json.Marshal(m)
	JSONPrint("", StructToString(m))
}

func checkFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// GetFilelist 获取指定目录的文件列表
func GetFilelist(path string) []string {
	files := make([]string, 0, 20)
	filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func getParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// GetViewHTML ...
func GetViewHTML(fileName string, data interface{}) string {
	var funcMaps = template.FuncMap{
		"empty": func(str string) bool {
			if str == "" {
				return true
			}
			return false
		},
		"unescaped": func(x string) interface{} { return template.HTML(x) },
	}
	strValue := strings.Split(fileName, "/")
	name := strValue[len(strValue)-1]

	pathName := fmt.Sprintf("%s", fileName)
	t, err := template.New(name).Funcs(funcMaps).ParseFiles(pathName)
	if err != nil {
		logger.Error(err)
	}

	html := bytes.NewBufferString("")
	err = t.Execute(html, data)
	if err != nil {
		log.Fatal(err)
		return "ERROR PARSE"
	}
	return html.String()
}

// GetLayoutViewHTML ...
func GetLayoutViewHTML(fileName string, data interface{}) string {
	var funcMaps = template.FuncMap{
		"empty": func(str string) bool {
			if str == "" {
				return true
			}
			return false
		},
		// "unescaped": func(x string) interface{} { return template.HTML(x) },
	}
	strValue := strings.Split(fileName, "/")
	name := strValue[len(strValue)-1]

	html := bytes.NewBufferString("")
	pathName := fmt.Sprintf("template/%s", fileName)
	tpl, err := template.New(name).Funcs(funcMaps).ParseFiles(pathName, "template/layout.html")
	if err != nil {
		logger.Error(err)
	}

	err = tpl.Execute(html, data)
	if err != nil {
		logger.Error(err)
	}

	return html.String()
}

// PathExists 判断路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// CopyFile 文件拷贝
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()

	os.Remove(dstName)
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}
