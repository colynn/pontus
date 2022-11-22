package sysadmin

import (
	"github.com/colynn/pontus/internal/account/role"
	"github.com/colynn/pontus/internal/menu"
	"github.com/colynn/pontus/internal/pkg/accountinfo"
	"github.com/colynn/pontus/internal/pkg/customerror"
	"github.com/colynn/pontus/internal/pkg/response"
	"github.com/colynn/pontus/tools"

	"github.com/gin-gonic/gin"
	log "unknwon.dev/clog/v2"
)

// GetMenuList ..
func GetMenuList(c *gin.Context) {
	Menu := &menu.Menu{}
	Menu.MenuName = c.Request.FormValue("menuName")
	Menu.Visible = c.Request.FormValue("visible")
	Menu.Title = c.Request.FormValue("title")
	// Menu.DataScope = tools.GetUserIdStr(c)
	result, err := Menu.GetMenuList()
	customerror.HasError(err, "抱歉未找到相关信息", 500)
	response.OK(c, result, "")
}

// GetMenuTreeelect ..
// For: add role
func GetMenuTreeelect(c *gin.Context) {
	var data menu.Menu
	result, err := data.SetMenuLable()
	customerror.HasError(err, "抱歉未找到相关信息", 500)
	response.OK(c, result, "")
}

// GetRoleMenu ..
func GetRoleMenu(c *gin.Context) {
	var Menu menu.Menu
	roleName := accountinfo.GetRoleName(c)
	log.Trace("get role menu, roleName: %v", roleName)
	result, err := Menu.GetMenuListByRole(roleName)
	customerror.HasError(err, "获取失败", 500)
	response.OK(c, result, "")
}

// GetMenuTreeRoleselect ..
func GetMenuTreeRoleselect(c *gin.Context) {
	var Menu menu.Menu
	var SysRole role.SysRole
	id, err := tools.StringToInt(c.Param("roleID"))
	customerror.HasError(err, "", 500)
	log.Trace("roleID: %d", id)
	SysRole.RoleID = id
	result, err := Menu.SetMenuLable()
	customerror.HasError(err, "抱歉未找到相关信息", 500)
	menuIds := make([]int, 0)
	if id != 0 {
		menuIds, err = SysRole.GetRoleMeunID()
		customerror.HasError(err, "抱歉未找到相关信息", 500)
	}
	response.Custom(c, gin.H{
		"code":        200,
		"menus":       result,
		"checkedKeys": menuIds,
	})
}

// GetMenu ..
func GetMenu(c *gin.Context) {
	var data menu.Menu
	id, err := tools.StringToInt(c.Param("id"))
	data.MenuID = id
	result, err := data.GetByMenuID()
	customerror.HasError(err, "抱歉未找到相关信息", 500)
	response.OK(c, result, "")
}

// InsertMenu ..
func InsertMenu(c *gin.Context) {
	var data menu.Menu
	err := c.ShouldBindJSON(&data)
	customerror.HasError(err, "抱歉未找到相关信息", 500)
	data.CreateBy = accountinfo.GetUserIDStr(c)
	result, err := data.Create()
	customerror.HasError(err, "抱歉未找到相关信息", 500)
	response.OK(c, result, "")
}

// UpdateMenu ..
func UpdateMenu(c *gin.Context) {
	var data menu.Menu
	err := c.ShouldBindJSON(&data)
	menuID, err := tools.StringToInt(c.Param("id"))
	data.UpdateBy = accountinfo.GetUserName(c)
	customerror.HasError(err, "修改失败", 500)
	_, err = data.Update(menuID)
	customerror.HasError(err, "", 501)
	response.OK(c, "", "修改成功")

}

// DeleteMenu ..
func DeleteMenu(c *gin.Context) {
	var data menu.Menu
	id, err := tools.StringToInt(c.Param("id"))
	data.UpdateBy = accountinfo.GetUserName(c)
	_, err = data.Delete(id)
	customerror.HasError(err, "删除失败", 500)
	response.OK(c, "", "删除成功")
}

// // GetMenuTreeRoleselect ..
// func GetMenuTreeRoleselect(c *gin.Context) {
// 	var Menu menu.Menu
// 	var SysRole role.SysRole
// 	id, err := tools.StringToInt(c.Param("roleId"))
// 	SysRole.RoleID = id
// 	result, err := Menu.SetMenuLable()
// 	customerror.HasError(err, "抱歉未找到相关信息", -1)
// 	MenuIDs := make([]int, 0)
// 	if id != 0 {
// 		MenuIDs, err = SysRole.GetRoleMeunId()
// 		customerror.HasError(err, "抱歉未找到相关信息", -1)
// 	}
// 	response.Custom(c, gin.H{
// 		"code":        200,
// 		"menus":       result,
// 		"checkedKeys": MenuIDs,
// 	})
// }

// // GetMenuIDS ..
// func GetMenuIDS(c *gin.Context) {
// 	var data models.RoleMenu
// 	data.RoleName = c.GetString("role")
// 	data.UpdateBy = tools.GetUserIdStr(c)
// 	result, err := data.GetIDS()
// 	customerror.HasError(err, "获取失败", 500)
// 	response.OK(c, result, "")
// }
