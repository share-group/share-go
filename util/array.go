package util

type arrayUtil struct{}

var ArrayUtil = newArrayUtil()

func newArrayUtil() *arrayUtil {
	return &arrayUtil{}
}

// 取数组最后一位元素
func (s *arrayUtil) Last(arr []string) string {
	if len(arr) <= 0 {
		return ""
	}
	return arr[len(arr)-1]
}
