package util

import (
	"fmt"
	"hash/fnv"
	"im/pkg/logger"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(UnixMilliTime(time.Now())))
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// RandString 生成随机字符串
func RandString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bytes[i] = byte(b)
	}
	return string(bytes)
}
func RandStringInt(num int) string {
	if num < 1 {
		return ""
	}

	letter := "0123456789abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, num)
	b[0] = letter[r.Intn(len(letter))]

	letter += "0"
	for i := 1; i < num; i++ {
		b[i] = letter[r.Intn(len(letter))]
	}
	return string(b)
}
func String2Int64(valS string) int64 {
	valInt64, err := strconv.ParseInt(valS, 10, 64)
	if err != nil {
		logger.Sugar.Error(GetSelfFuncName(), "convert string to int64 failed, err: %s", err)
	}
	return valInt64
}
func String2Int(valS string) int {
	valInt, err := strconv.Atoi(valS)
	if err != nil {
		logger.Sugar.Error(GetSelfFuncName(), "convert string to int failed, err: %s", err)
	}
	return valInt
}

func Int2RuneFromString(valS string) rune {
	return rune(String2Int(valS))
}

func IndexOf(list []string, value string) int {
	for k, v := range list {
		if v == value {
			return k
		}
	}
	return -1
}

// InSliceForString 判断字符串是否在 slice 中。
func InSliceForString(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// InSliceForInt64 判断int64类型是否在 slice 中。
func InSliceForInt64(items []int64, item int64) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func Uint2Sting(u uint) string {
	return strconv.Itoa(int(u))
}

func Int64Sting(u int64) string {
	return strconv.Itoa(int(u))
}

func KeyMatch(key1 string, key2 string) bool {
	key2 = strings.Replace(key2, "/*", "/.*", -1)

	re := regexp.MustCompile(`:[^/]+`)
	key2 = re.ReplaceAllString(key2, "$1[^/]+$2")

	return RegexMatch(key1, "^"+key2+"$")
}

func RegexMatch(key1 string, key2 string) bool {
	res, err := regexp.MatchString(key2, key1)
	if err != nil {
		panic(err)
	}
	return res
}

func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func StringToInt(i string) int {
	j, _ := strconv.Atoi(i)
	return j
}
func StringToInt64(i string) int64 {
	j, _ := strconv.ParseInt(i, 10, 64)
	return j
}
func StringToInt32(i string) int32 {
	j, _ := strconv.ParseInt(i, 10, 64)
	return int32(j)
}
func Int32ToString(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func MatchPhone(str string) bool {
	reg := regexp.MustCompile(`\+?\d{11}`) // 正则调用规则，regex包不支持?=之类的格式，后续有需求再换regex2
	return reg.MatchString(str)            // 返回 MatchString 是否匹配
}

func RandInt(min, max int64) int64 {
	return min + r.Int63n(max-min)
}

func RandID(num int) string {
	if num < 1 {
		return ""
	}

	letter := "123456789"
	b := make([]byte, num)
	b[0] = letter[r.Intn(len(letter))]

	letter += "0"
	for i := 1; i < num; i++ {
		b[i] = letter[r.Intn(len(letter))]
	}
	return string(b)
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func RandLimitInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

func StringHash(str string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(str))
	return h.Sum32()
}

// Int64SliceToString 下述方法里边的3行代码，任意一行都可以实现将int类型切片转化成指定格式的字符串
func Int64SliceToString(a []int64, delim string) string {
	//return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
	//return strings.Trim(strings.Join(strings.Split(fmt.Sprint(a), " "), delim), "[]")
	return strings.Trim(strings.Replace(fmt.Sprint(a), " ", delim, -1), "[]")
}
