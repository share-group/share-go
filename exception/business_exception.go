package errors

import (
	"fmt"
)

// BusinessException 是一个实现了 error 接口的结构体
type BusinessException struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 创建一个新的异常实例
func NewBusinessException(code int, message string) BusinessException {
	return BusinessException{Code: code, Message: message}
}

// 抛出异常
func Throw(businessException BusinessException) {
	panic(businessException)
}

// 实现 error 接口的 Error 方法
func (b *BusinessException) Error() string {
	return fmt.Sprintf("BusinessException: Code %d - %s", b.Code, b.Message)
}
