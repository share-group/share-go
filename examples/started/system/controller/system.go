package system

import (
	"github.com/share-group/share-go/bootstrap"
	system "github.com/share-group/share-go/examples/started/system/protocol"
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

func (c *systemController) GetCaptcha() *system.ResponseCaptcha {
	if len(c.Api) > 0 {
		//fmt.Println(1 / (len(c.Api) - 1))
		//return nil, errors.New("这个是个错误")
	}
	log.Info("哈哈哈  哦哦哦   嘻嘻😳😳")
	return &system.ResponseCaptcha{UUID: "xxxxxx", Captcha: "xxxxxxxxxx"}
}
