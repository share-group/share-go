package util

import (
	"math/rand"
)

type randomUtil struct{}

var RandomUtil = newRandomUtil()

func newRandomUtil() *randomUtil {
	return &randomUtil{}
}

// 生成随机字符串
//
// length-字符串长度
func (r *randomUtil) String(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Int(0, len(bytes)-1)])
	}
	return string(result)
}

// 生成随机整数
//
// min-最小值
// max-最大值
func (r *randomUtil) Int(min int, max int) int {
	return rand.Intn(max-min+1) + min
}
