package system

import (
	"github.com/share-group/share-go/examples/started/protocol"
	"github.com/share-group/share-go/examples/started/service"
)

type captchaController struct{}

var CaptchaController = newCaptchaController()

func newCaptchaController() *captchaController {
	return &captchaController{}
}

func (c *captchaController) GetCaptcha() *protocol.ResponseCaptcha {
	id, b64s := service.CaptchaService.GetCaptcha()
	return &protocol.ResponseCaptcha{UUID: id, Captcha: b64s}
}
