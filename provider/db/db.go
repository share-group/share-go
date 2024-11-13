package db

import (
	"github.com/share-group/share-go/provider/db/mongodb"
)

// 注册数据表实体
//
// entity-数据实体; connectionNames-连接名称
func RegisterRepository(entity any, connectionNames ...string) {
	for _, connectionName := range connectionNames {
		InjectEntity(entity, connectionName).EnsureIndex(entity)
	}
}

// 获取数据表实体
//
// entity-数据实体; connectionName-连接名称
func InjectEntity[T any](entity T, connectionName ...string) *mongodb.Mongodb[T] {
	_connectionName := "default"
	if len(connectionName) > 0 {
		_connectionName = connectionName[0]
	}
	return mongodb.GetInstance(entity, _connectionName)
}
