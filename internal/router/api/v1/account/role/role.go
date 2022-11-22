package role

import (
	"github.com/colynn/pontus/internal/menu"
	"github.com/colynn/pontus/internal/pkg/customerror"
	"github.com/colynn/pontus/internal/pkg/response"
	"github.com/colynn/pontus/tools"

	"github.com/colynn/pontus/internal/account/role"

	"github.com/gin-gonic/gin"
	log "unknwon.dev/clog/v2"
)

// GetSysRoleList ..
func GetSysRoleList(c *gin.Context) {
	var data role.SysRole
	var err error
	var pageSize = 10
	var pageIndex = 1
	size := c.Request.FormValue("pageSize")
	if size != "" {
		pageSize = tools.StrToInt(err, size)
	}

	index := c.Request.FormValue("pageIndex")
	if index != "" {
		pageIndex = tools.StrToInt(err, index)
	}

	data.RoleKey = c.Request.FormValue("roleKey")
	data.RoleName = c.Request.FormValue("roleName")
	data.Status = c.Request.FormValue("status")
	log.Trace("request params: pageIndex: %v, pageSize: %v", pageIndex, pageSize)
	result, count, err := data.GetPage(pageSize, pageIndex)
	customerror.HasError(err, "", -1)

	response.PageOK(c, result, count, pageIndex, pageSize, "")
}

// GetAllRoles ..
func GetAllRoles(c *gin.Context) {
	var data role.SysRole
	result, err := data.GetList()
	customerror.HasError(err, "", 500)
	response.OK(c, result, "success")
}

// GetRole ..
func GetRole(c *gin.Context) {
	data := role.SysRole{}
	data.RoleID, _ = tools.StringToInt(c.Param("roleID"))
	result, err := data.Get()
	menuIds := make([]int, 0)
	menuIds, err = data.GetRoleMeunID()
	customerror.HasError(err, "抱歉未找到相关信息", 500)
	result.MenuIds = menuIds
	response.OK(c, result, "")
}

// CreateRole ..
func CreateRole(c *gin.Context) {
	data := role.SysRole{}
	c.ShouldBindJSON(&data)
	ID, err := data.Insert()
	customerror.HasError(err, "", 500)

	rm := menu.RoleMenu{}
	rm.RoleID = ID
	_, err = rm.Insert(ID, data.MenuIds)
	customerror.HasError(err, "", 500)

	// result.MenuIds = menuIds
	response.OK(c, "success", "")
}

// UpdateRole ..
func UpdateRole(c *gin.Context) {
	data := role.SysRole{}
	data.RoleID, _ = tools.StringToInt(c.Param("roleID"))
	c.ShouldBindJSON(&data)
	result, err := data.Update(data.RoleID)

	var t menu.RoleMenu
	_, err = t.DeleteRoleMenu(data.RoleID)
	customerror.HasError(err, "添加失败1", 500)
	_, err2 := t.Insert(data.RoleID, data.MenuIds)
	customerror.HasError(err2, "添加失败2", 500)

	// result.MenuIds = menuIds
	response.OK(c, result, "")
}

// DeleteRole ..
func DeleteRole(c *gin.Context) {
	data := role.SysRole{}
	RoleID, _ := tools.StringToInt(c.Param("roleID"))
	ids := []int{RoleID}
	result, err := data.BatchDelete(ids)
	// TODO: menuIDs
	// menuIds := make([]int, 0)
	// menuIds, err = data.GetRoleMeunId()
	customerror.HasError(err, "抱歉未找到相关信息", -1)

	// result.MenuIds = menuIds
	response.OK(c, result, "")
}
