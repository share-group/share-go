package system

import (
	UserError "github.com/share-group/share-go/examples/started/error/user"
	"github.com/share-group/share-go/examples/started/protocol"
	"github.com/share-group/share-go/examples/started/service"
	"github.com/share-group/share-go/provider/logger"
	"github.com/share-group/share-go/provider/mongodb"
	"github.com/share-group/share-go/provider/redis"
	"github.com/share-group/share-go/util"
	"strconv"
)

type captchaController struct{}

var CaptchaController = newCaptchaController()

func newCaptchaController() *captchaController {
	return &captchaController{}
}

type RequestCaptcha struct {
	A string `json:"a" validate:"required" message:"aaaaaa"`
	B string `json:"b" validate:"required" message:"bbbbbbbbb"`
}

func (c *captchaController) GetCaptcha(r *RequestCaptcha) *protocol.ResponseCaptcha {
	id, b64s := service.CaptchaService.GetCaptcha()
	logger.GetLogger().Info("xxxxxxxxxxxxxx")
	util.SystemUtil.AssertAndThrowError(r.B == "1", UserError.TokenError)

	if r.B == "2" {
		b, _ := strconv.Atoi(r.B)
		logger.GetLogger().Info(strconv.Itoa(2 / (b - 2)))
	}
	return &protocol.ResponseCaptcha{UUID: id, Captcha: b64s}
}
func (c *captchaController) Captcha2(r *protocol.RequestLogin) *protocol.ResponseCaptcha {
	id, b64s := service.CaptchaService.GetCaptcha()
	redis.Redis.Client.Get("aa")
	mongodb.Mongodb.DB.Name()
	return &protocol.ResponseCaptcha{UUID: id, Captcha: b64s}
}
