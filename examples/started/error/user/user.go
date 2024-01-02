package UserError

import exception "github.com/share-group/share-go/exception"

// 20001-系统开了个小差，请稍后重试
var TokenError = exception.NewBusinessException(20001, "登录信息已失效，请重新登录")
