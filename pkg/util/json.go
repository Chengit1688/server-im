package util

import (
	"encoding/json"
	"github.com/json-iterator/go"
	"reflect"
	"strconv"
)

func JsonMarshal(v interface{}) ([]byte, error) {
	return jsoniter.Marshal(v)
}

func JsonUnmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}

// Copy src结构体内容拷贝到dst
func Copy(src, dst interface{}) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, dst)
	if err != nil {
		return err
	}
	return nil
}

func MarshalJSONByDefault(i interface{}, returnJson bool) (interface{}, error) {
	typeOf := reflect.TypeOf(i)
	valueOf := reflect.ValueOf(i)
	for i := 0; i < typeOf.Elem().NumField(); i++ {
		if valueOf.Elem().Field(i).IsZero() {
			def := typeOf.Elem().Field(i).Tag.Get("default")
			zo := typeOf.Elem().Field(i).Tag.Get("zero")
			if def != "" {
				switch typeOf.Elem().Field(i).Type.String() {
				case "int64":
					if zo != "" {
						valueOf.Elem().Field(i).SetInt(0)
					} else {
						result := String2Int64(def)
						valueOf.Elem().Field(i).SetInt(result)
					}
				case "uint":
					if zo != "" {
						valueOf.Elem().Field(i).SetUint(0)
					} else {
						result, _ := strconv.ParseUint(def, 10, 64)
						valueOf.Elem().Field(i).SetUint(result)
					}
				case "string":
					if zo == "" {
						valueOf.Elem().Field(i).SetString("")
					} else {
						valueOf.Elem().Field(i).SetString(def)
					}
				}
			}
		}
	}
	if !returnJson {
		return i, nil
	}
	return json.Marshal(i)
}
