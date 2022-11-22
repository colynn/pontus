package initialize

import (
	"github.com/colynn/pontus/internal/account/role"

	log "unknwon.dev/clog/v2"
)

// init role
// admin / common

// InitRole ensure have admin/common user
func InitRole() error {
	//
	roleModel := role.SysRole{}
	// TODO: role_id did not AUTO_INCREMENT correctly
	initRoleList := []role.SysRole{
		{
			RoleKey:  "admin",
			RoleName: "超级管理员",
			Status:   "0",
		},
		{
			RoleKey:  "common",
			RoleName: "普通用户",
			Status:   "0",
		},
		{
			RoleKey:  "developer",
			RoleName: "开发者",
			Status:   "0",
		},
		{
			RoleKey:  "default",
			RoleName: "默认用户",
			Status:   "0",
		},
	}
	for _, item := range initRoleList {
		roleModel.RoleKey = item.RoleKey
		_, err := roleModel.Get()
		if err == nil {
			// rolekey already exits
			continue
		}
		log.Warn("when initialize role, %s role not exists, error: %s, recreate it", item.RoleKey, err.Error())
		roleModel.RoleName = item.RoleName
		roleModel.Status = item.Status
		roleID, err := roleModel.Insert()
		if err != nil {
			log.Error("when initialize role %s, insert error: %s", item.RoleKey, err.Error())
			continue
		}
		log.Info("initialize role %s success, role id: %v", item.RoleKey, roleID)
	}
	return nil
}
