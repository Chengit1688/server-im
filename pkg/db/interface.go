package db

import (
	"errors"
	"fmt"
	"im/internal/api/group/model"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

/*
map 查询下的高级支持
map下的value支持整数与文本，与整数数组 整数表示相等查询，可查0 整数数组表示 in 操作该数组
文本指令已特殊字符开头
>参数 大于参数 <参数 小于参数
>=参数 <=参数
^参数1,参数2 取两者之间
?参数 模糊查询
=参数 等于该参数
无特殊指令则为相等判断

/缺少 or 支持 缺少不等于支持

*/

// 通用性单条数据查询函数
// 传入 model
// 传入where条件 支持 传入id，id列表,struct结构,map 结构
// 传入id时取id对应数据，struct有效数据作为where条件,map 支持部分高级查询，具体点击此处查看
func Info(mod interface{}, where interface{}) (err error) {
	typeMod, _, _ := GetInterfaceRef(mod)
	if typeMod.Kind() != reflect.Struct {
		return errors.New("model 参数错误，不是struct结构")
	}
	query := DB.Model(mod)
	query = CheckWhere(query, typeMod, where)

	err = query.First(mod).Error

	return err
}

func InfoTx(tx *gorm.DB, mod interface{}, where interface{}) (err error) {
	typeMod, _, _ := GetInterfaceRef(mod)
	if typeMod.Kind() != reflect.Struct {
		return errors.New("model 参数错误，不是struct结构")
	}
	query := tx.Model(mod)
	query = CheckWhere(query, typeMod, where)

	err = query.First(mod).Error

	return err
}

// 获取指定字段的输出列表 如跨表时获取群id列表等等
func CloumnList(mod interface{}, where interface{}, cloumn string) (out []interface{}, err error) {
	typeMod, _, _ := GetInterfaceRef(mod)
	if typeMod.Kind() != reflect.Struct {
		return nil, errors.New("model 参数错误，不是struct结构")
	}
	query := DB.Model(mod)
	query = CheckWhere(query, typeMod, where)
	data := [](map[string]interface{}){}
	err = query.Select(cloumn).Find(&data).Error
	if err != nil {
		return
	}
	out = []interface{}{}
	for _, v := range data {
		value, had := v[cloumn]
		if !had {
			continue
		}
		out = append(out, value)
	}
	return out, err
}

// 通用性列表查询函数
// 传入 model
// 传入where条件 支持 传入id，id列表,struct结构,map 结构
// 传入排序参数 如 'id desc' 'id asc'
// 分页数
// 每页数据
// 指针返回 总数
// 指针返回 结果列表
func Find(mod interface{}, where interface{}, sort string, page, pageSize int, total *int64, out interface{}) (err error) {
	typeMod, _, _ := GetInterfaceRef(mod)
	if typeMod.Kind() != reflect.Struct {
		return errors.New("model 参数错误，不是struct结构")
	}
	query := DB.Debug().Model(mod)
	query = CheckWhere(query, typeMod, where)
	query.Count(total)
	if sort != "" {
		query = query.Order(sort)
	}
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	err = query.Find(out).Error
	return err
}

func Count(mod interface{}, where interface{}) (total int64) {
	typeMod, _, _ := GetInterfaceRef(mod)
	if typeMod.Kind() != reflect.Struct {
		return
	}
	query := DB.Debug().Model(mod)
	query = CheckWhere(query, typeMod, where)
	query.Count(&total)

	return
}

func Insert(data interface{}) error {
	err := DB.Create(data).Error
	return err
}

func Update(mod interface{}, where interface{}, data interface{}) error {
	query := DB.Model(mod)
	typeMod, _, _ := GetInterfaceRef(mod)
	query = CheckWhere(query, typeMod, where)
	err := query.Updates(data).Error
	return err
}

func Delete(mod interface{}, where interface{}) error {
	query := DB.Model(mod)
	typeMod, _, _ := GetInterfaceRef(mod)
	query = CheckWhere(query, typeMod, where)
	err := query.Unscoped().Delete(mod).Error
	return err
}

func InsertTx(tx *gorm.DB, data interface{}) error {
	err := tx.Create(data).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

func UpdateTx(tx *gorm.DB, mod interface{}, where interface{}, data interface{}) error {
	query := tx.Model(mod).Debug()
	typeMod, _, _ := GetInterfaceRef(mod)
	query = CheckWhere(query, typeMod, where)
	err := query.Updates(data).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

func DeleteTx(tx *gorm.DB, mod interface{}, where interface{}) error {
	query := tx.Model(mod)
	typeMod, _, _ := GetInterfaceRef(mod)
	query = CheckWhere(query, typeMod, where)
	err := query.Unscoped().Delete(mod).Error
	if err != nil {
		tx.Rollback()
	}
	return err
}

func FindTx(tx *gorm.DB, mod interface{}, where interface{}, sort string, page, pageSize int, total *int64, out interface{}) (err error) {
	typeMod, _, _ := GetInterfaceRef(mod)
	if typeMod.Kind() != reflect.Struct {
		return errors.New("model 参数错误，不是struct结构")
	}
	query := tx.Debug().Model(mod)
	query = CheckWhere(query, typeMod, where)
	query.Count(total)
	if sort != "" {
		query = query.Order(sort)
	}
	if page > 0 && pageSize > 0 {
		offset := (page - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}
	err = query.Find(out).Error
	return err
}
func GetInterfaceRef(mod interface{}) (reflect.Type, reflect.Value, bool) {
	typeOfA := reflect.TypeOf(mod)
	valueOfA := reflect.ValueOf(mod)
	isPtr := false
	if typeOfA.Kind() == reflect.Ptr {
		typeOfA = typeOfA.Elem()
		valueOfA = valueOfA.Elem()
		isPtr = true
	}
	return typeOfA, valueOfA, isPtr
}

func CheckWhere(query *gorm.DB, mod reflect.Type, where interface{}) *gorm.DB {
	if where == nil {
		return query
	}
	typeData, valueData, _ := GetInterfaceRef(where)
	switch typeData.Kind() {
	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		//寻找主键id

		priKey := GetPrimary(mod)

		if priKey == "" {
			return query
		}
		query = query.Where(fmt.Sprintf("`%s` = ?", priKey), valueData.Int())
		break
	case reflect.String:
		//寻找主键id
		priKey := GetPrimary(mod)
		if priKey == "" {
			return query
		}
		query = query.Where(fmt.Sprintf("`%s` = ?", priKey), valueData.String())
		break
	case reflect.Struct:
		whereMap := GetStructMap(typeData, valueData)
		for k, v := range whereMap {
			query = query.Where(fmt.Sprintf("`%s` = ?", k), v)
		}
		break
	case reflect.Map:
		//锁定map[string]interface结构
		// whereMap := where.(map[string]interface{})
		query = WhereMapTrans(query, mod, where.(map[string]interface{}))
		break

	}
	return query
}

func GetPrimary(mod reflect.Type) string {
	hadId := false
	for i := 0; i < mod.NumField(); i++ {
		tagValue := mod.Field(i).Tag.Get("gorm")
		tagValue = strings.ToLower(tagValue)
		tagValueMap := GetTagKV(tagValue)
		if strings.ToLower(mod.Field(i).Name) == "id" {
			hadId = true
		}
		if _, had := tagValueMap["primary_key"]; had {
			//返回主键
			if _, hadName := tagValueMap["column"]; hadName {
				return tagValueMap["column"]
			}
			return mod.Field(i).Name
		}
	}

	if hadId {
		return "id"
	}
	return ""

}

func GetColumnName(mod reflect.StructField) string {
	tagValue := mod.Tag.Get("gorm")
	tagValue = strings.ToLower(tagValue)
	tagValueMap := GetTagKV(tagValue)

	if _, hadName := tagValueMap["column"]; hadName {
		return tagValueMap["column"]
	}
	return mod.Name

}

func GetJsonName(mod reflect.StructField) string {
	tagValue := mod.Tag.Get("json")
	tagValue = strings.ToLower(tagValue)
	tagValueMap := GetTagKV(tagValue)

	for k, _ := range tagValueMap {
		return k
	}
	return mod.Name

}

// 将struct转换为有效map结构
func GetStructMap(typeMod reflect.Type, valueMod reflect.Value) map[string]interface{} {
	out := map[string]interface{}{}
	for i := 0; i < typeMod.NumField(); i++ {
		itemValue := valueMod.Field(i)
		if itemValue.IsValid() && !itemValue.IsZero() { //这里的判断函数可能有问题
			columnTag := GetColumnName(typeMod.Field(i))
			out[columnTag] = itemValue.Interface()
		}
	}
	return out
}

// func MakeOrWhere(keys []string,value )(key string,params []interface{}){

// }

// 直接输入的高级map 允许部分特殊操作
func WhereMapTrans(query *gorm.DB, modType reflect.Type, whereMap map[string]interface{}) *gorm.DB {
	//先建立所有json级 key与类型的记录
	kt := map[string]reflect.Kind{}
	for i := 0; i < modType.NumField(); i++ {
		itemValue := modType.Field(i)
		columnName := GetColumnName(itemValue)
		kt[columnName] = itemValue.Type.Kind()
	}
	for key, value := range whereMap {
		//寻找对应key的元素类型
		itemType, had := kt[key]
		keys := strings.Split(key, "|")
		if len(keys) == 1 && !had {
			continue
		} else {
			fmt.Println("or模式", key, keys)
			clearnKeys := []string{}
			for _, v := range keys {
				if itemType, had = kt[v]; had {
					clearnKeys = append(clearnKeys, v)
				}
			}
			keys = clearnKeys
			fmt.Println("or模式清理", keys)
			if len(keys) == 0 {
				continue
			}
		}

		switch itemType {
		case reflect.Int, reflect.Int32, reflect.Int16, reflect.Int64:
			//目标类型为int，则统一转换为int64进行查询
			switch value.(type) {
			case int, int16, int32, int64:
				if len(keys) == 1 {
					query = query.Where(fmt.Sprintf("`%s`=?", key), value)
				} else {
					keyStr := []string{}
					endInter := []interface{}{}
					for i := 0; i < len(keys); i++ {
						keyStr = append(keyStr, fmt.Sprintf("`%s`=?", keys[i]))
						endInter = append(endInter, value)
					}
					query = query.Where(fmt.Sprintf("(%s)", keyStr), endInter...)
				}

				break
			case string:
				action, int1, int2 := WhereStringParseToInt64(value.(string))
				switch action {
				case "", "=":
					if len(keys) == 1 {
						query = query.Where(fmt.Sprintf("`%s`=?", key), int1)
					} else {
						keyStr := []string{}
						endInter := []interface{}{}
						for i := 0; i < len(keys); i++ {
							keyStr = append(keyStr, fmt.Sprintf("`%s`=?", keys[i]))
							endInter = append(endInter, value)
						}
						query = query.Where(fmt.Sprintf("(%s)", strings.Join(keyStr, " or ")), endInter...)
					}
					break
				case ">", ">=", "<", "<=":

					if len(keys) == 1 {
						query = query.Where(fmt.Sprintf("`%s`%s?", key, action), int1)
					} else {
						keyStr := []string{}
						endInter := []interface{}{}
						for i := 0; i < len(keys); i++ {
							keyStr = append(keyStr, fmt.Sprintf("`%s`%s?", keys[i], action))
							endInter = append(endInter, int1)
						}
						query = query.Where(fmt.Sprintf("(%s)", strings.Join(keyStr, " or ")), endInter...)
					}
					break
				case "^":

					if len(keys) == 1 {
						query = query.Where(fmt.Sprintf("`%s` between ? and ?", key), int1, int2)
					} else {
						keyStr := []string{}
						endInter := []interface{}{}
						for i := 0; i < len(keys); i++ {
							keyStr = append(keyStr, fmt.Sprintf("`%s` between ? and ?", keys[i]))
							endInter = append(endInter, int1, int2)
						}
						query = query.Where(fmt.Sprintf("(%s)", strings.Join(keyStr, " or ")), endInter...)
					}
					break
				}
			default:
				//特殊类型里去支持数组
				if reflect.TypeOf(value).Kind() == reflect.Array || reflect.TypeOf(value).Kind() == reflect.Slice {

					if len(keys) == 1 {
						query = query.Where(fmt.Sprintf("`%s` in (?)", key), value)
					} else {
						keyStr := []string{}
						endInter := []interface{}{}
						for i := 0; i < len(keys); i++ {
							keyStr = append(keyStr, fmt.Sprintf("`%s` in (?)", keys[i]))
							endInter = append(endInter, value)
						}
						query = query.Where(fmt.Sprintf("(%s)", strings.Join(keyStr, " or ")), endInter...)
					}
				}
				// case array:
			}
		case reflect.String:

			//目标为字符串 则只允许字符串类型的特殊查询
			switch value.(type) {
			case model.RoleType:
				strVal := value.(model.RoleType)
				action, acValue := WhereStringParseToString(string(strVal))
				switch action {
				case "", "=":

					if len(keys) == 1 {
						query = query.Where(fmt.Sprintf("`%s`=?", key), acValue)
					} else {
						keyStr := []string{}
						endInter := []interface{}{}
						for i := 0; i < len(keys); i++ {
							keyStr = append(keyStr, fmt.Sprintf("`%s`=?", keys[i]))
							endInter = append(endInter, acValue)
						}
						query = query.Where(fmt.Sprintf("(%s)", strings.Join(keyStr, " or ")), endInter...)
					}
				}
			case string:
				action, acValue := WhereStringParseToString(value.(string))

				switch action {
				case "", "=":

					if len(keys) == 1 {
						query = query.Where(fmt.Sprintf("`%s`=?", key), acValue)
					} else {
						keyStr := []string{}
						endInter := []interface{}{}
						for i := 0; i < len(keys); i++ {
							keyStr = append(keyStr, fmt.Sprintf("`%s`=?", keys[i]))
							endInter = append(endInter, acValue)
						}
						query = query.Where(fmt.Sprintf("(%s)", strings.Join(keyStr, " or ")), endInter...)
					}
					break
				case ">", ">=", "<", "<=":

					if len(keys) == 1 {
						query = query.Where(fmt.Sprintf("`%s`%s?", key, action), acValue)
					} else {
						keyStr := []string{}
						endInter := []interface{}{}
						for i := 0; i < len(keys); i++ {
							keyStr = append(keyStr, fmt.Sprintf("`%s`%s?", keys[i], action))
							endInter = append(endInter, acValue)
						}
						query = query.Where(fmt.Sprintf("(%s)", strings.Join(keyStr, " or ")), endInter...)
					}
				case "?":

					if len(keys) == 1 {
						query = query.Where(fmt.Sprintf("`%s` like ?", key), "%"+acValue+"%")
					} else {
						keyStr := []string{}
						endInter := []interface{}{}
						for i := 0; i < len(keys); i++ {
							keyStr = append(keyStr, fmt.Sprintf("`%s` like ?", keys[i]))
							endInter = append(endInter, "%"+acValue+"%")
						}
						query = query.Where(fmt.Sprintf("(%s)", strings.Join(keyStr, " or ")), endInter...)
					}
				}

			default:
				//特殊类型里去支持数组
				if reflect.TypeOf(value).Kind() == reflect.Array || reflect.TypeOf(value).Kind() == reflect.Slice {
					//
					if len(keys) == 1 {
						query = query.Where(fmt.Sprintf("`%s` in (?)", key), value)
					} else {
						keyStr := []string{}
						endInter := []interface{}{}
						for i := 0; i < len(keys); i++ {
							keyStr = append(keyStr, fmt.Sprintf("`%s` in (?)", keys[i]))
							endInter = append(endInter, value)
						}
						query = query.Where(fmt.Sprintf("(%s)", strings.Join(keyStr, " or ")), endInter...)
					}
				}
				// case array:
			}
		}
	}
	return query
}

func MapToStruct(data map[string]interface{}, typeMod reflect.Type) (interface{}, error) {
	//先建立所有column级 key与类型的记录
	kt := map[string]int{}
	for i := 0; i < typeMod.NumField(); i++ {
		itemValue := typeMod.Field(i)
		columnName := GetColumnName(itemValue)
		kt[columnName] = i
	}
	newData := reflect.New(typeMod).Elem()
	for k, v := range data {
		if index, had := kt[k]; had {
			switch newData.Field(index).Kind() {
			case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
				if v == nil {
					newData.Field(index).SetInt(0)
				} else {
					newData.Field(index).SetInt(reflect.ValueOf(v).Int())
				}

			case reflect.String:
				newData.Field(index).SetString(reflect.ValueOf(v).String())
			case reflect.Float32, reflect.Float64:
				newData.Field(index).SetFloat(reflect.ValueOf(v).Float())
			}

		}
	}

	return newData.Interface(), nil
}

// 将where字符串解析为条件与值
func WhereStringParseToString(str string) (string, string) {
	reg := regexp.MustCompile(`^(<=|>=|\?|\^|>|<|=)`)
	action := reg.FindString(str)
	value := strings.Replace(str, action, "", 1)
	value = strings.TrimSpace(value)
	return action, value
}

// 将where字符串解析为条件与数字参数
func WhereStringParseToInt64(str string) (string, int64, int64) {
	reg := regexp.MustCompile(`^(<=|>=|\^|>|<|=)`)
	action := reg.FindString(str)
	value := strings.Replace(str, action, "", 1)
	if action == "^" {
		//区间查询，特殊处理
		nums := strings.Split(value, ",")
		if len(nums) >= 2 {
			nums[0] = strings.TrimSpace(nums[0])
			nums[1] = strings.TrimSpace(nums[1])
			start, _ := strconv.ParseInt(nums[0], 10, 64)
			end, _ := strconv.ParseInt(nums[1], 10, 64)
			return action, start, end
		}
	}
	value = strings.TrimSpace(value)
	valueInt, _ := strconv.ParseInt(value, 10, 64)
	return action, valueInt, 0
}

func GetTagKV(str string) map[string]string {
	lists := strings.Split(str, ";")
	out := map[string]string{}
	for i := 0; i < len(lists); i++ {
		kvs := strings.Split(lists[i], ":")
		if len(kvs) >= 2 {
			out[kvs[0]] = kvs[1]
		}

		if len(kvs) == 1 {
			out[kvs[0]] = ""
		}
	}

	return out
}
