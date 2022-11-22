package gorm

import (
	"github.com/colynn/pontus/config"

	"github.com/colynn/pontus/internal/account/role"
	"github.com/colynn/pontus/internal/account/user"
	"github.com/colynn/pontus/internal/audit"
	"github.com/colynn/pontus/internal/cmdb"
	"github.com/colynn/pontus/internal/menu"
	mycasbin "github.com/colynn/pontus/internal/middleware/casbin"

	"github.com/jinzhu/gorm"
)

// AutoMigrate ..
func AutoMigrate(db *gorm.DB) error {
	c := config.GetConfig()
	if c.GetString("database.dbtype") == "mysql" {
		db = db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	}
	db.SingularTable(true)
	return db.AutoMigrate(
		new(mycasbin.CasbinRule),
		new(user.SysUser),
		new(user.LoginLog),
		new(role.SysRole),
		new(cmdb.CloudAsset),
		new(cmdb.AssetTag),
		new(cmdb.TangibleAsset),
		new(cmdb.AssetUser),
		new(audit.SysAudit),
		new(menu.Menu),
		new(menu.RoleMenu),
	).Error
}
