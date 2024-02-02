package system

import (
	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/examples/started/protocol"
	"github.com/share-group/share-go/provider/mongodb"
	"github.com/share-group/share-go/provider/redis"
)

type captchaController struct{}

var CaptchaController = newCaptchaController()

var mongo = mongodb.GetInstance("dashboard")
var r = redis.GetInstance()

func newCaptchaController() *captchaController {
	return &captchaController{}
}

func (o *captchaController) GetCaptcha(c echo.Context) *protocol.ResponseCaptcha {
	return &protocol.ResponseCaptcha{}
}
