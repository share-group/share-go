package system

import (
	"github.com/share-group/share-go/bootstrap"
)

var log = bootstrap.Logger.GetLogger()

type systemController struct {
	Api    string
	Method string
}

var SystemController = newSystemController()

func newSystemController() *systemController {
	return &systemController{
		Api:    "a",
		Method: "GET",
	}
}

func (c *systemController) GetCaptcha() string {
	log.Info("哈哈哈  哦哦哦   嘻嘻😳😳")
	//return &system.ResponseCaptcha{UUID: "xxxxxx", Captcha: "xxxxxxxxxx"}
	return "1212"
}

func (c *systemController) Login() string {
	log.Info("哈哈哈  哦哦哦   嘻嘻😳😳")
	//return &system.ResponseCaptcha{UUID: "xxxxxx", Captcha: "xxxxxxxxxx"}
	return "1212"
}
