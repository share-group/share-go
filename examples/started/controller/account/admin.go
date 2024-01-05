package account

import (
	"github.com/share-group/share-go/examples/started/protocol"
	"log"
)

type adminController struct{}

var AdminController = newAdminController()

func newAdminController() *adminController {
	return &adminController{}
}

func (a *adminController) Login(r *protocol.RequestLogin) int {
	log.Println(r)
	return 1
}
