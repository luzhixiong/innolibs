package utils

import (
	"fmt"
	"reflect"
	"strings"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: StructToMap
//@description: 利用反射将结构体转化为map
//@param: obj interface{}
//@return: map[string]interface{}

func StructToMap(obj interface{}) map[string]interface{} {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	data := make(map[string]interface{})
	for i := 0; i < obj1.NumField(); i++ {
		if obj1.Field(i).Tag.Get("mapstructure") != "" {
			data[obj1.Field(i).Tag.Get("mapstructure")] = obj2.Field(i).Interface()
		} else {
			data[obj1.Field(i).Name] = obj2.Field(i).Interface()
		}
	}
	return data
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: ArrayToString
//@description: 将数组格式化为字符串
//@param: array []interface{}
//@return: string

func ArrayToString(array []interface{}) string {
	return strings.Replace(strings.Trim(fmt.Sprint(array), "[]"), " ", ",", -1)
}

// 判断字符是不是在数组内
func InArray(needle string, haystack []string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}
func InArrayUint(needle uint, haystack []uint) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

// 判斷
func InArrayStruct(needle struct {
	Key   string
	Value int32
	Sort  int
}, haystack []struct {
	Key   string
	Value int32
	Sort  int
}) bool {
	for _, item := range haystack {
		if item.Key == needle.Key && item.Value == needle.Value && item.Sort == needle.Sort {
			return true
		}
	}

	return false
}
