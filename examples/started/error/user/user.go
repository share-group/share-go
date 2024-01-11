package UserError

import exception "github.com/share-group/share-go/exception"

// 20001-登录信息已失效，请重新登录
var TokenError = exception.NewBusinessException(20001, "登录信息已失效，请重新登录")
