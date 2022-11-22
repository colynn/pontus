package user

import (
	"errors"
	"net/http"

	"github.com/colynn/pontus/internal/account/role"
	"github.com/colynn/pontus/internal/menu"
	"github.com/colynn/pontus/internal/pkg/accountinfo"
	userinfo "github.com/colynn/pontus/internal/pkg/accountinfo"
	"github.com/colynn/pontus/internal/pkg/customerror"
	"github.com/colynn/pontus/internal/pkg/response"
	"github.com/colynn/pontus/tools"
	"github.com/colynn/pontus/tools/currenttime"
	"github.com/colynn/pontus/tools/file"
	"github.com/colynn/pontus/tools/locationip"

	"github.com/colynn/pontus/tools/captcha"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	log "unknwon.dev/clog/v2"
)

// GenerateCaptchaHandler ..
func GenerateCaptchaHandler(c *gin.Context) {
	id, b64s, err := captcha.DriverDigitFunc()
	customerror.HasError(err, "验证码获取失败", 500)
	response.Custom(c, gin.H{
		"code": 200,
		"data": b64s,
		"id":   id,
		"msg":  "success",
	})
}

// GetSysUserList ..
func GetSysUserList(c *gin.Context) {
	var data SysUser
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

	data.Username = c.Request.FormValue("username")

	result, count, err := data.GetPage(pageSize, pageIndex)
	customerror.HasError(err, "", -1)

	response.PageOK(c, result, count, pageIndex, pageSize, "")
}

// GetSysUser ..
func GetSysUser(c *gin.Context) {
	var SysUser SysUser
	UserID, _ := tools.StringToInt(c.Param("userID"))
	result, err := SysUser.Get(UserID)
	customerror.HasError(err, "抱歉未找到相关信息", 500)
	var SysRole role.SysRole
	roles, err := SysRole.GetList()

	roleIds := make([]int, 0, 1)
	roleIds = append(roleIds, result.RoleID)
	response.Custom(c, gin.H{
		"code":    200,
		"data":    result,
		"roleIds": roleIds,
		"roles":   roles,
	})
}

// CreateUser ..
func CreateUser(c *gin.Context) {
	userModel := SysUser{}
	c.ShouldBindJSON(&userModel)
	result, err := userModel.Insert(&userModel)
	customerror.HasError(err, "", 500)
	response.OK(c, result, "")
}

// UpdateUser ..
func UpdateUser(c *gin.Context) {
	userModel := SysUser{}
	userID, err := tools.StringToInt(c.Param("userID"))
	if err != nil {
		log.Error("get user id occur error: %s", err.Error())
	}
	log.Trace("userID: %d", userID)
	userModel.SysUserId.ID = userID
	c.ShouldBindJSON(&userModel)
	result, err := userModel.Update(userID)
	customerror.HasError(err, "抱歉未找到相关信息", 500)
	response.OK(c, result, "")
}

// UpdateUserProfile ..
func UpdateUserProfile(c *gin.Context) {
	userModel := SysUser{}
	userID := accountinfo.GetUserID(c)
	log.Trace("userID: %d", userID)
	userModel.SysUserId.ID = userID
	c.ShouldBindJSON(&userModel)
	result, err := userModel.Update(userID)
	customerror.HasError(err, "抱歉未找到相关信息", 500)
	response.OK(c, result, "")
}

// ImportUserData ..
func ImportUserData(c *gin.Context) {
	f, fh, err := c.Request.FormFile("file")
	customerror.HasError(err, "", 500)

	if fh.Size > 1024*1024 {
		customerror.HasError(errors.New("导入文件不允许超过1M"), "", 500)
	}
	defer f.Close()
	filePath, err := file.StorageImportFile(f, fh)
	customerror.HasError(err, "", 500)
	userSvc := NewService()
	err = userSvc.InsertOrUpdateUserData(filePath)
	customerror.HasError(err, "", 500)
	response.OK(c, filePath, "添加成功")
}

// DeleteSysUser ..
func DeleteSysUser(c *gin.Context) {
	var data SysUser
	userID, err := tools.StringToInt(c.Param("userID"))
	customerror.HasError(err, "", 500)
	ids := []int{userID}
	result, err := data.BatchDelete(ids)
	customerror.HasError(err, "删除失败", 500)
	response.OK(c, result, "删除成功")
}

// GetInfo ..
func GetInfo(c *gin.Context) {

	sysuser := SysUser{}
	userID := userinfo.GetUserID(c)
	if userID == 0 {
		err := errors.New("未获取到用户信息，请重新登录")
		customerror.HasError(err, "", 500)
	}
	user, err := sysuser.Get(userID)
	customerror.HasError(err, "", 500)

	var mp = make(map[string]interface{})

	mp["avatar"] = "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif"
	if user.Avatar != "" {
		mp["avatar"] = user.Avatar
	}
	mp["realName"] = user.RealName
	mp["username"] = user.Username
	mp["email"] = user.Email
	mp["phone"] = user.Phone
	mp["id"] = user.ID

	AdminPermissions := "*:*:*"
	mp["roles"] = []string{userinfo.GetRoleName(c)}

	if userinfo.GetRoleName(c) == "admin" {
		mp["permissions"] = []string{AdminPermissions}
	} else {
		roleMenu := menu.RoleMenu{}
		roleMenu.RoleID = userinfo.GetRoleID(c)
		log.Trace("role id: %d", roleMenu.RoleID)
		list, _ := roleMenu.GetPermis()
		mp["permissions"] = list
	}
	response.OK(c, mp, "")
}

// LogOut ..
func LogOut(c *gin.Context) {
	var loginlog LoginLog
	ua := user_agent.New(c.Request.UserAgent())
	loginlog.Ipaddr = c.ClientIP()
	location := locationip.GetLocation(c.ClientIP())
	loginlog.LoginLocation = location
	loginlog.LoginTime = currenttime.GetCurrentTime()
	loginlog.Status = "0"
	loginlog.Remark = c.Request.UserAgent()
	browserName, browserVersion := ua.Browser()
	loginlog.Browser = browserName + " " + browserVersion
	loginlog.Os = ua.OS()
	loginlog.Platform = ua.Platform()
	loginlog.Username = userinfo.GetUserName(c)
	loginlog.Msg = "退出成功"
	loginlog.Create()
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "退出成功",
	})

}
