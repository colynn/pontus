package menu

import (
	"fmt"

	"github.com/colynn/pontus/internal/account/role"
	"github.com/colynn/pontus/internal/db"
	"github.com/colynn/pontus/tools"

	log "unknwon.dev/clog/v2"
)

// RoleMenu ..
type RoleMenu struct {
	RoleID   int    `gorm:"type:int(11)"`
	MenuID   int    `gorm:"type:int(11)"`
	RoleName string `gorm:"type:varchar(128)"`
	CreateBy string `gorm:"type:varchar(128)"`
	UpdateBy string `gorm:"type:varchar(128)"`
}

// TableName ..
func (RoleMenu) TableName() string {
	return "sys_role_menu"
}

// Path ..
type Path struct {
	Path string `json:"path"`
}

// Get ..
func (rm *RoleMenu) Get() ([]RoleMenu, error) {
	var r []RoleMenu
	table := db.Eloquent.Table("sys_role_menu")
	if rm.RoleID != 0 {
		table = table.Where("role_id = ?", rm.RoleID)

	}
	if err := table.Find(&r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

// GetPermis ..
func (rm *RoleMenu) GetPermis() ([]string, error) {
	var r []Menu
	table := db.Eloquent.Select("sys_menu.permission").Table("sys_menu").Joins("left join sys_role_menu on sys_menu.menu_id = sys_role_menu.menu_id")

	table = table.Where("role_id = ?", rm.RoleID)

	table = table.Where("sys_menu.menu_type in (?, ?)", Button, MenuType)
	if err := table.Find(&r).Error; err != nil {
		return nil, err
	}
	list := make([]string, 0, len(r))
	for i := 0; i < len(r); i++ {
		list = append(list, r[i].Permission)
	}
	return list, nil
}

// GetIDS ..
func (rm *RoleMenu) GetIDS() ([]int, error) {
	var r []Menu
	// TODO: change to tableName
	table := db.Eloquent.Select("sys_menu.menu_id").Table("sys_menu").Where("menu_type in (?, ?, ?)", Directory, MenuType, APIInterface)
	if err := table.Find(&r).Error; err != nil {
		return nil, err
	}
	list := make([]int, 0, len(r))
	for i := 0; i < len(r); i++ {
		list = append(list, r[i].MenuID)
	}
	return list, nil
}

// GetPaths ..
func (rm *RoleMenu) GetPaths() ([]Path, error) {
	var r []Path
	table := db.Eloquent.Select("sys_menu.path").Table("sys_role_menu")
	table = table.Joins("left join sys_role on sys_role.role_id=sys_role_menu.role_id")
	table = table.Joins("left join sys_menu on sys_menu.id=sys_role_menu.menu_id")
	table = table.Where("sys_role.role_name = ?", rm.RoleName)
	if err := table.Find(&r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

// DeleteRoleMenu ..
func (rm *RoleMenu) DeleteRoleMenu(roleID int) (bool, error) {
	if err := db.Eloquent.Table("sys_role_menu").Where("role_id = ?", roleID).Delete(&rm).Error; err != nil {
		return false, err
	}
	var roleModel role.SysRole
	if err := db.Eloquent.Table("sys_role").Where("role_id = ?", roleID).First(&roleModel).Error; err != nil {
		return false, err
	}
	sql3 := "delete from casbin_rule where v0= '" + roleModel.RoleKey + "';"
	db.Eloquent.Exec(sql3)

	return true, nil

}

// BatchDeleteRoleMenu ..
func (rm *RoleMenu) BatchDeleteRoleMenu(roleIds []int) (bool, error) {
	if err := db.Eloquent.Table("sys_role_menu").Where("role_id in (?)", roleIds).Delete(&rm).Error; err != nil {
		return false, err
	}
	var roleItems []role.SysRole
	if err := db.Eloquent.Table("sys_role").Where("role_id in (?)", roleIds).Find(&roleItems).Error; err != nil {
		return false, err
	}
	sql := ""
	for i := 0; i < len(roleItems); i++ {
		sql += "delete from casbin_rule where v0= '" + roleItems[i].RoleName + "';"
	}
	db.Eloquent.Exec(sql)
	return true, nil

}

// Insert ..
func (rm *RoleMenu) Insert(roleID int, menuID []int) (bool, error) {
	var roleModel role.SysRole
	if err := db.Eloquent.Table("sys_role").Where("role_id = ?", roleID).First(&roleModel).Error; err != nil {
		log.Error("get role by id: %v error: %s", roleID, err.Error())
		return false, err
	}
	var menu []Menu
	if err := db.Eloquent.Table("sys_menu").Where("menu_id in (?)", menuID).Find(&menu).Error; err != nil {
		log.Error("get menus by ids: %v error: %s", menuID, err.Error())
		return false, err
	}
	//ORM不支持批量插入所以需要拼接 sql 串
	sql := "INSERT INTO `sys_role_menu` (`role_id`,`menu_id`,`role_name`) VALUES "
	originSQLLength := len(sql)
	sql2 := "INSERT INTO casbin_rule  (`p_type`,`v0`,`v1`,`v2`) VALUES "
	originLength := len(sql2)
	for i := 0; i < len(menu); i++ {
		if len(menu)-1 == i {
			//最后一条数据 以分号结尾
			sql += fmt.Sprintf("(%d,%d,'%s');", roleModel.RoleID, menu[i].MenuID, roleModel.RoleKey)
			if menu[i].MenuType == "A" {
				sql2 += fmt.Sprintf("('p','%s','%s','%s');", roleModel.RoleKey, menu[i].Path, menu[i].Action)
			}
		} else {
			sql += fmt.Sprintf("(%d,%d,'%s'),", roleModel.RoleID, menu[i].MenuID, roleModel.RoleKey)
			if menu[i].MenuType == "A" {
				sql2 += fmt.Sprintf("('p','%s','%s','%s'),", roleModel.RoleKey, menu[i].Path, menu[i].Action)
			}
		}
	}

	log.Trace("role menu sql: %s\n", sql)
	if originSQLLength == len(sql) {
		log.Warn("role menu sql is invalidation, skip exec")
	} else {
		db.Eloquent.Exec(sql)
	}
	casbinSQL := sql2[0:len(sql2)-1] + ";"
	log.Trace("sql2: %s, %d", sql2, originLength)
	log.Trace("casbin sql: %s, %d", casbinSQL, len(casbinSQL))
	if originLength == len(casbinSQL) {
		log.Warn("casbin sql is invalidation, skip exec")
	} else {
		db.Eloquent.Exec(casbinSQL)
	}

	return true, nil
}

// Delete ..
func (rm *RoleMenu) Delete(RoleID string, MenuID string) (bool, error) {
	rm.RoleID, _ = tools.StringToInt(RoleID)
	table := db.Eloquent.Table("sys_role_menu").Where("role_id = ?", RoleID)
	if MenuID != "" {
		table = table.Where("menu_id = ?", MenuID)
	}
	if err := table.Delete(&rm).Error; err != nil {
		return false, err
	}
	return true, nil

}
