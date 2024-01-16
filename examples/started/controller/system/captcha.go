package system

import (
	"github.com/share-group/share-go/examples/started/protocol"
	"github.com/share-group/share-go/examples/started/service"
	"github.com/share-group/share-go/provider/logger"
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

func (c *captchaController) GetCaptcha() *protocol.ResponseCaptcha {
	id, b64s := service.CaptchaService.GetCaptcha()
	return &protocol.ResponseCaptcha{UUID: id, Captcha: b64s}
}
func (c *captchaController) Captcha2() {
	logger.GetLogger().Info("1111111111111111111111111111111")
}
