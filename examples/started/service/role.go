package service

import (
	"fmt"
	entity "github.com/share-group/share-go/examples/started/entity/account"
	"github.com/share-group/share-go/provider/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

type roleService struct {
}

var RoleService = newRoleService()

func newRoleService() *roleService {
	return &roleService{}
}

var mongo = mongodb.GetInstance("dashboard")

func init() {
	mongo.UpdateOne(entity.Role{}, "65af4b58bf294f8c9745c7d7", bson.D{{"index", 431}})
	roleList := RoleService.RoleList("", 1, 1)
	if len(roleList) > 0 {
		return
	}
	RoleService.createRole("超级管理员")
}

// 角色列表
//
// name-角色名称搜索;page-当前页码;pageSize-页面大小
func (r *roleService) RoleList(name string, page, pageSize int64) []entity.Role {
	query := make(bson.D, 0)
	if len(name) > 0 {
		query = append(query, bson.E{Key: "name", Value: name})
	}

	ctx, cursor := mongo.PaginationByPage(query, page, pageSize, make(bson.D, 0), entity.Role{})
	return mongodb.Decode(ctx, cursor, entity.Role{})
}

// 创建角色
//
// name-角色名称
func (r *roleService) createRole(name string) {
	fmt.Println("id: ", mongo.InsertMany(entity.Role{Name: name, Authority: make([]any, 0)}, entity.Role{Name: "aaaa", Authority: make([]any, 0)}, entity.Role{Name: "bbb", Authority: make([]any, 0)}))
}
