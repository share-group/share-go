package system

import (
	"github.com/share-group/share-go/examples/started/protocol"
	"github.com/share-group/share-go/examples/started/service"
	"log"
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
	return &protocol.ResponseCaptcha{UUID: id, Captcha: b64s}
}
func (c *captchaController) Captcha2(r *protocol.RequestLogin) *protocol.ResponseCaptcha {
	id, b64s := service.CaptchaService.GetCaptcha()
	log.Println("1111 ", r)
	return &protocol.ResponseCaptcha{UUID: id, Captcha: b64s}
}
