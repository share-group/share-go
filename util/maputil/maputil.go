package maputil

import (
	"fmt"
	"reflect"
	"strings"
)

// GetValueFromMap 从映射中获取值，如果键不存在则返回默认值
func GetValueFromMap[K comparable, V any](m map[K]V, key K, defaultValue V) V {
	if m == nil {
		return defaultValue
	}
	if val, ok := m[key]; ok {
		if reflect.DeepEqual("<nil>", fmt.Sprintf("%v", val)) || reflect.ValueOf(val).IsZero() {
			return defaultValue
		}
		return val
	}
	return defaultValue
}

// ContainsKey 判断映射中是否包含指定键
func ContainsKey[K comparable, V any](m map[K]V, key K) bool {
	if _, ok := m[key]; ok {
		return true
	}
	return false
}

// Keys 返回map所有键
func Keys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	return keys
}

// Merge 多个map合并
func Merge[K comparable, V any](maps ...map[K]V) map[K]V {
	target := make(map[K]V)
	for _, m := range maps {
		for k, v := range m {
			target[k] = v
		}
	}
	return target
}

// kv翻转
func Flip[K comparable, V comparable](m map[K]V) map[V]K {
	tmp := make(map[V]K)
	for k, v := range m {
		tmp[v] = k
	}
	return tmp
}

// kv转url
func Map2url[K comparable, V comparable](data map[K]V) string {
	list := make([]string, 0)
	for k, v := range data {
		list = append(list, strings.TrimSpace(fmt.Sprintf("%v=%v", k, v)))
	}
	return strings.Join(list, "&")
}
