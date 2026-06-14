package account

import (
	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/examples/started/protocol"
	"github.com/share-group/share-go/examples/started/util/echoutil"
	loggerFactory "github.com/share-group/share-go/provider/logger"
)

type roleController struct{}

var logger = loggerFactory.GetLogger()
var RoleController = newRoleController()

func newRoleController() *roleController {
	logger.Info("newRoleController")
	echoutil.Echo()
	return &roleController{}
}

func (a *roleController) Login(c echo.Context, r *protocol.RequestLogin) int {
	logger.Info("%v  %v", r, c.RealIP())
	return 1
}
