package system

import "github.com/share-group/share-go/bootstrap"

type logService struct{}

var LogService = newLogService()

var logMongodb = bootstrap.NewMongodb(bootstrap.Config.GetStringValue("data.logging.uri"))

func newLogService() *logService {
	return &logService{}
}

func (s *logService) SaveRequestLog() {
}
