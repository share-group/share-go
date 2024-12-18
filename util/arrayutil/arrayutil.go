package arrayutil

import (
	"reflect"
)

// 取数组第一位元素
//
// arr-数组
func First[T any](arr []T) (t T) {
	if len(arr) <= 0 {
		return t
	}
	return arr[0]
}

// 取数组最后一位元素
//
// arr-数组
func Last[T any](arr []T) (t T) {
	if len(arr) <= 0 {
		return t
	}
	return arr[len(arr)-1]
}

// 判断元素是否存在于数组中
func Contains[T any](arr []T, element T) bool {
	for _, item := range arr {
		if reflect.DeepEqual(item, element) {
			return true
		}
	}
	return false
}

// 对象数组元素去重
func Uniq[T comparable](array []T) []T {
	newArray := make([]T, 0)
	encountered := map[T]bool{}

	for _, v := range array {
		if encountered[v] != true {
			encountered[v] = true
			newArray = append(newArray, v)
		}
	}
	return newArray
}

// 接受一个过滤条件，返回一个新数组
func Filter[T any](arr []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range arr {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}
