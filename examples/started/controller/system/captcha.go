package system

import (
	"fmt"
	entity "github.com/share-group/share-go/examples/started/entity/account"
	"github.com/share-group/share-go/examples/started/protocol"
	"github.com/share-group/share-go/examples/started/service"
	"github.com/share-group/share-go/provider/mongodb"
	"github.com/share-group/share-go/provider/redis"
	"time"
)

type captchaController struct{}

var CaptchaController = newCaptchaController()

var mongo = mongodb.GetInstance("dashboard")
var r = redis.GetInstance()

func newCaptchaController() *captchaController {
	return &captchaController{}
}

func (c *captchaController) GetCaptcha() *protocol.ResponseCaptcha {
	id, b64s := service.CaptchaService.GetCaptcha()
	r.Client.SetNX("1", 1, 3*time.Minute)
	fmt.Println(mongo.DB.Name(), entity.Admin{})
	return &protocol.ResponseCaptcha{UUID: id, Captcha: b64s}
}
