package util

type arrayUtil struct{}

var ArrayUtil = newArrayUtil()

func newArrayUtil() *arrayUtil {
	return &arrayUtil{}
}

// 取数组第一位元素
//
// arr-数组
func (s *arrayUtil) First(arr []string) string {
	if len(arr) <= 0 {
		return ""
	}
	return arr[0]
}

// 取数组最后一位元素
//
// arr-数组
func (s *arrayUtil) Last(arr []string) string {
	if len(arr) <= 0 {
		return ""
	}
	return arr[len(arr)-1]
}
