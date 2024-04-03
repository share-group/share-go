package account

import (
	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/examples/started/protocol"
	"log"
)

type roleController struct{}

var RoleController = newRoleController()

func newRoleController() *roleController {
	return &roleController{}
}

func (a *roleController) Login(c echo.Context, r *protocol.RequestLogin) int {
	log.Println(r, c.RealIP())
	return 1
}
