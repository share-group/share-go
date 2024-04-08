package service

import (
	entity "github.com/share-group/share-go/examples/started/entity/tubemax"
	"github.com/share-group/share-go/provider/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

var tubemax = mongodb.GetInstance("tubemax")

type systemService struct{}

var SystemService = newSystemService()

func newSystemService() *systemService {
	return &systemService{}
}

func (s *systemService) GetGlobalConfig() []*entity.SystemConfig {
	ctx, cursor := tubemax.Find(bson.D{}, entity.SystemConfig{})
	return mongodb.DecodeList(ctx, cursor, entity.SystemConfig{})
}
