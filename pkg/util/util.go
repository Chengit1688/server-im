package util

import (
	"fmt"
	"math/rand"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func GetSelfFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return cleanUpFuncName(runtime.FuncForPC(pc).Name())
}

func cleanUpFuncName(funcName string) string {
	end := strings.LastIndex(funcName, ".")
	if end == -1 {
		return ""
	}
	return funcName[end+1:]
}

func GetFriendResKey(kind string, name string) string {
	//schedule:groups:group111
	pre := GetFriendResKeyPrefix(kind)
	return fmt.Sprintf("%s%s", pre, name)
}

func GetFriendResKeyPrefix(kind string) string {
	return fmt.Sprintf("friend:%s:", kind)
}
func Intersect(slice1, slice2 []int) []int {
	m := make(map[int]bool)
	n := make([]int, 0)
	for _, v := range slice1 {
		m[v] = true
	}
	for _, v := range slice2 {
		flag, _ := m[v]
		if flag {
			n = append(n, v)
		}
	}
	return n
}

// 获取 slice1 比 slice2 多的元素
func Difference(slice1, slice2 []int) []int {
	m := make(map[int]bool)
	n := make([]int, 0)
	inter := Intersect(slice1, slice2)
	for _, v := range inter {
		m[v] = true
	}
	for _, v := range slice1 {
		if !m[v] {
			n = append(n, v)
		}
	}
	// 注释 暂时不需要slice2 比slice1多的
	// for _, v := range slice2 {
	// 	if !m[v] {
	// 		n = append(n, v)
	// 	}
	// }
	return n
}

// 只允许数字和字母
func IsAlphaNumeric(s string) bool {
	regex := regexp.MustCompile("^[a-zA-Z0-9]+$")
	return regex.MatchString(s)
}

// 只允许数字、字母和中文
func IsAlphaNumericChinese(s string) bool {
	regex := regexp.MustCompile("^[\\p{Han}\\p{Latin}\\p{Nd}]*$")
	return regex.MatchString(s)
}

// 提前定义能抢到的最小金额1分
var min int64 = 1

// 二倍均值算法
func DoubleAverage(count, amount int64) int64 {
	if count == 1 {
		return amount
	}
	//计算出最大可用金额
	max := amount - min*count
	//计算出最大可用平均值
	avg := max / count
	//二倍均值基础上再加上最小金额 防止出现金额为0
	avg2 := 2*avg + min
	//随机红包金额序列元素，把二倍均值作为随机的最大数
	rand.Seed(time.Now().UnixNano())
	x := rand.Int63n(avg2) + min
	return x
}

func RemoveDuplicateElement(addrs []string) []string {
	result := make([]string, 0, len(addrs))
	temp := map[string]struct{}{}
	for _, item := range addrs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
