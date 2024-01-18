package arrayutil

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
