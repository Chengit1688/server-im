package util

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jinzhu/copier"
)

func CopyStructFields(dst interface{}, src interface{}, fields ...string) (err error) {
	return copier.Copy(dst, src)
}

func StructToMap(in interface{}, tagName string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			tagValueList := strings.Split(tagValue, ",")
			if len(tagValueList) > 1 {
				tagValue = tagValueList[0]
			}
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil
}

func StructToMapWithoutNil(in interface{}, tagName string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			tagValueList := strings.Split(tagValue, ",")
			if len(tagValueList) > 1 {
				tagValue = tagValueList[0]
			}

			fieldValue := v.Field(i)
			if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
				out[tagValue] = fieldValue.Elem().Interface()
			} else if fieldValue.Kind() != reflect.Ptr {
				out[tagValue] = fieldValue.Interface()
			}
		}
	}

	return out, nil
}
