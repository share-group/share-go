package account

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/examples/started/protocol"
	loggerFactory "github.com/share-group/share-go/provider/logger"
)

type roleController struct{}

var RoleController = newRoleController()

func newRoleController() *roleController {
	loggerFactory.GetLogger().Info("xxxxxxxxxxxx")
	return &roleController{}
}

func (a *roleController) Login(c echo.Context, r *protocol.RequestLogin) int {
	log.Println(r, c.RealIP())
	return 1
}
