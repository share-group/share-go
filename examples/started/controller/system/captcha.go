package system

import (
	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/examples/started/protocol"
	"log"
)

type captchaController struct{}

var CaptchaController = newCaptchaController()

func newCaptchaController() *captchaController {
	return &captchaController{}
}

func (o *captchaController) GetCaptcha(c echo.Context) *protocol.ResponseCaptcha {
	return &protocol.ResponseCaptcha{}
}

func (o *captchaController) GetNumber(c echo.Context, r *protocol.RequestLogin) int {
	log.Println(r)
	log.Println(c.Request().Header.Get("Userid"))
	log.Println(c.Request().Header.Get("userId"))
	return 1
}
