package SystemError

import (
	exception "github.com/share-group/share-go/exception"
)

// 10001-系统开了个小差，请稍后重试
var SystemError = exception.NewBusinessException(10001, "系统开了个小差，请稍后重试")

// 10002-公共参数错误
var CommonParametersError = exception.NewBusinessException(10002, "公共参数错误")
