package system

import (
	"fmt"
	"github.com/labstack/gommon/random"
	system "github.com/share-group/share-go/examples/started/system/protocol"
	"go.uber.org/zap"
	"time"
)

type systemController struct{}

var SystemController = newSystemController()

func newSystemController() *systemController {
	return &systemController{}
}

func (c *systemController) Login(request *system.LoginRequest) *system.LoginResponse {
	zap.L().Info(fmt.Sprintf("request: %v", request))
	uuid := fmt.Sprintf("%s_%d", random.String(32), time.Now().UnixMilli())
	return &system.LoginResponse{Token: uuid}
}
