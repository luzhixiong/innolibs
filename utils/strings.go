package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oklog/ulid/v2"
	"math"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	char           = "abcdefghijklmnopqrstuvwxyz"
	charCapital    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	number         = "0123456789"
	charNum        = "abcdefghijklmnopqrstuvwxyz0123456789"
	charCapitalNum = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charMix        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

	rentSalt = "flyebike2023"
)

// 字符串首个字符大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// 字符串首个字符小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// MaheHump 将字符串转换为驼峰命名
func MaheHump(s string) string {
	words := strings.Split(s, "-")

	for i := 1; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}

	return strings.Join(words, "")
}

// 随机生成字符串
func randStr(size int, str string) string {
	var s bytes.Buffer
	for i := 0; i < size; i++ {
		s.WriteByte(str[Rander().Int63()%int64(len(str))])
	}
	return s.String()
}

func RandChar(size int) string {
	return randStr(size, char)
}

func RandCapitalChar(size int) string {
	return randStr(size, charCapital)
}

func RandNum(size int) string {
	return randStr(size, number)
}

func RandCharNum(size int) string {
	return randStr(size, charNum)
}

func RandHex(size int) string {
	b := make([]byte, size)
	Rander().Read(b)
	return fmt.Sprintf("%x", b)[:size]
}

func RandCapitalCharNum(size int) string {
	return randStr(size, charCapitalNum)
}

func RandMix(size int) string {
	return randStr(size, charMix)
}

// 字符串补0
func StrPad(input string, padLength int, padString string, padType string) string {
	output := ""
	inputLen := len(input)

	if inputLen >= padLength {
		return input
	}

	padStringLen := len(padString)
	needFillLen := padLength - inputLen

	if diffLen := padStringLen - needFillLen; diffLen > 0 {
		padString = padString[diffLen:]
	}

	for i := 1; i <= needFillLen; i += padStringLen {
		output += padString
	}
	switch padType {
	case "LEFT":
		return output + input
	default:
		return input + output
	}
}

func HEX2DEC(hex string) int64 {
	dec, _ := strconv.ParseInt(hex, 16, 64)
	return dec
}

func IF(isTrue bool, a, b interface{}) interface{} {
	if isTrue {
		return a
	}
	return b
}

// 字符串中包含的字符[]
func RemoveStrs(str string, subs []string) string {
	for _, s := range subs {
		str = RemoveStr(str, s)
	}
	return str
}

// 字符串中包含的字符
func RemoveStr(str, sub string) string {
	return strings.ReplaceAll(str, sub, "")
}

// 加*
func HideMiddleStr(str string, first, last int) string {
	length := len(str)
	if first+last < length {
		tmp := str[:first]
		lenT := length - first - last
		for i := 0; i < lenT; i++ {
			tmp += "*"
		}
		tmp += str[length-last:]
		return tmp
	} else {
		return str
	}
}

// 对比版本号
func CompareVersion(version1 string, version2 string) int {
	var res int
	ver1Strs := strings.Split(version1, ".")
	ver2Strs := strings.Split(version2, ".")
	ver1Len := len(ver1Strs)
	ver2Len := len(ver2Strs)
	verLen := ver1Len
	if len(ver1Strs) < len(ver2Strs) {
		verLen = ver2Len
	}
	for i := 0; i < verLen; i++ {
		var ver1Int, ver2Int int
		if i < ver1Len {
			ver1Int, _ = strconv.Atoi(ver1Strs[i])
		}
		if i < ver2Len {
			ver2Int, _ = strconv.Atoi(ver2Strs[i])
		}
		if ver1Int < ver2Int {
			res = -1
			break
		}
		if ver1Int > ver2Int {
			res = 1
			break
		}
	}
	return res
}

func UnescapeHTML(content string) string {
	content = strings.ReplaceAll(content, "&quot;", `"`)
	content = strings.ReplaceAll(content, "&apos;", `'`)
	content = strings.ReplaceAll(content, "&lt;", `<`)
	content = strings.ReplaceAll(content, "&gt;", `>`)
	content = strings.ReplaceAll(content, "&ldquo;", `「`)
	content = strings.ReplaceAll(content, "&rdquo;", `」`)
	content = strings.ReplaceAll(content, "&amp;", `&`) // 最后处理
	return content
}

func RentMd5Salt(password string) string {
	return MD5V([]byte(password + rentSalt))
}

// 生成26位随机字符串，替换uuid
func ULID() string {
	ms := ulid.Timestamp(time.Now())
	id, _ := ulid.New(ms, Rander())
	return strings.ToLower(id.String())
}

// Any2String 任意类型转为字符串
func Any2String(any interface{}) string {
	if any == nil {
		return ""
	}
	switch value := any.(type) {
	case int:
		return strconv.Itoa(value)
	case int8:
		return strconv.Itoa(int(value))
	case int16:
		return strconv.Itoa(int(value))
	case int32:
		return strconv.Itoa(int(value))
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint8:
		return strconv.FormatUint(uint64(value), 10)
	case uint16:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(value)
	case string:
		return value
	case []byte:
		return string(value)
	case time.Time:
		if value.IsZero() {
			return ""
		}
		return value.String()
	case *time.Time:
		if value == nil {
			return ""
		}
		return value.String()
	default:
		if value == nil {
			return ""
		}
		var (
			rv   = reflect.ValueOf(value)
			kind = rv.Kind()
		)
		switch kind {
		case reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return ""
			}
		case reflect.String:
			return rv.String()
		}
		if kind == reflect.Ptr {
			return Any2String(rv.Elem().Interface())
		}
		if jsonContent, err := json.Marshal(value); err != nil {
			return fmt.Sprint(value)
		} else {
			return string(jsonContent)
		}
	}
}

func GetMsgIds(now time.Time) (traceId, messageId string) {
	traceId = fmt.Sprintf("0%s%d%s", RandHex(7), now.UnixNano()/100, RandHex(5))
	messageId = fmt.Sprintf("%d", now.UnixNano())
	return
}

func CorrectJsonStr(str string) string {
	placer := strings.NewReplacer(
		"，", ",",
		" ", "",
		"\r\n", "",
		"\r", "",
		"\n", "",
		"\t", "",
	)
	str = placer.Replace(str)
	length := len(str)
	if length < 3 {
		return str
	}

	if str[length-2:length-1] == "," {
		str = str[:length-2] + str[length-1:]
	}
	return str
}

func Ipv4ToLong(ip string) (uint, error) {
	p := net.ParseIP(ip).To4()
	if p == nil {
		return 0, errors.New("invalid ipv4 format")
	}

	return uint(p[0])<<24 | uint(p[1])<<16 | uint(p[2])<<8 | uint(p[3]), nil
}

func LongToIpv4(i uint) (string, error) {
	if i > math.MaxUint32 {
		return "", errors.New("beyond the scope of ipv4")
	}

	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(i >> 24)
	ip[1] = byte(i >> 16)
	ip[2] = byte(i >> 8)
	ip[3] = byte(i)

	return ip.String(), nil
}
