package account

import (
	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/examples/started/protocol"
	"github.com/share-group/share-go/examples/started/service"
)

type adminController struct{}

var AdminController = newAdminController()

func newAdminController() *adminController {
	return &adminController{}
}

func (a *adminController) Login(c echo.Context, r *protocol.RequestLogin) any {
	return service.SystemService.GetGlobalConfig()
}
