package system

import (
	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/examples/started/protocol"
)

type captchaController struct{}

var CaptchaController = newCaptchaController()

func newCaptchaController() *captchaController {
	return &captchaController{}
}

func (o *captchaController) GetCaptcha(c echo.Context) *protocol.ResponseCaptcha {
	return &protocol.ResponseCaptcha{}
}
