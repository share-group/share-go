package account

import (
	"github.com/labstack/echo/v4"
	"github.com/share-group/share-go/examples/started/protocol"
	"log"
)

type adminController struct{}

var AdminController = newAdminController()

func newAdminController() *adminController {
	return &adminController{}
}

func (a *adminController) Login(c echo.Context, r *protocol.RequestLogin) int {
	log.Println(r, c.RealIP())
	return 1
}
