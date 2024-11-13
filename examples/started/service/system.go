package service

import (
	entity "github.com/share-group/share-go/examples/started/entity/tubemax"
	"github.com/share-group/share-go/provider/db/mongodb"
)

var tubemax = mongodb.GetInstance(entity.SystemConfig{}, "tubemax")

type systemService struct{}

var SystemService = newSystemService()

func newSystemService() *systemService {
	return &systemService{}
}

func (s *systemService) GetGlobalConfig() []*entity.SystemConfig {
	//ctx, cursor := tubemax.Find(bson.D{}, entity.SystemConfig{})
	return make([]*entity.SystemConfig, 0) // mongodb.DecodeList(ctx, cursor, entity.SystemConfig{})
}
