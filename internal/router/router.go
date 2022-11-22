package router

import (
	"fmt"
	"log"

	"github.com/colynn/pontus/internal/middleware"

	"github.com/colynn/pontus/internal/account/user"
	"github.com/colynn/pontus/internal/cmdb"
	"github.com/colynn/pontus/internal/router/api/v1/account/role"
	"github.com/colynn/pontus/internal/router/api/v1/audit"
	"github.com/colynn/pontus/internal/router/api/v1/sysadmin"

	apiv1cmdb "github.com/colynn/pontus/internal/router/api/v1/cmdb"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// InitRouter ..
func InitRouter(cfg *viper.Viper) *gin.Engine {
	r := gin.New()
	r.Use(middleware.LoggerToFile())
	r.Use(middleware.HandleError)
	r.Use(middleware.NoCache)
	r.Use(middleware.Options)
	r.Use(middleware.Secure)

	// the jwt middleware
	authMiddleware, err := middleware.AuthInit()

	if err != nil {
		_ = fmt.Errorf("JWT Error: %v", err.Error())
	}

	r.POST("/login", authMiddleware.LoginHandler)

	// Refresh time can be longer than token timeout
	r.GET("/refresh_token", authMiddleware.RefreshHandler)

	apiv1 := r.Group("/api/v1")
	{
		apiv1.GET("/getCaptcha", user.GenerateCaptchaHandler)
		apiv1.POST("/tags", cmdb.CreateTag)
		apiv1.POST("/assets/clouds", cmdb.CreateMutilpleServer)
		apiv1.GET("/assets/statistics", cmdb.AssetStatistics)
		apiv1.GET("/assets/clouds/export", cmdb.ExportCloudAsset)
		apiv1.POST("/logout", user.LogOut)
	}

	auth := r.Group("/api/v1")
	// auth.Use(authMiddleware.MiddlewareFunc())
	auth.Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())

	{

		// 资产管理
		// auth.GET("/assets/clouds/export", cmdb.ExportCloudAsset)
		auth.GET("/assets/clouds", cmdb.GetServersByPagination)
		auth.GET("/assets/pc/list", apiv1cmdb.GetPClist)
		auth.POST("/assets/pc/upload", apiv1cmdb.ImportTerminalDeviceData)
		// physical
		auth.POST("/assets/physical/upload", apiv1cmdb.ImportPhysicalData)
		auth.GET("/assets/physical/list", apiv1cmdb.GetPhysicallist)

		auth.GET("/assets/tangibles/:instanceID", apiv1cmdb.GetTangibleAssetInfo)
		auth.PUT("/assets/tangibles/:instanceID", middleware.SysAudit(), apiv1cmdb.UpdateTangibleAssetInfo)
		auth.DELETE("/assets/tangibles/:instanceID", middleware.SysAudit(), apiv1cmdb.DeleteTangibleAssetInfo)
		auth.GET("/assets/tangibles/:instanceID/audits", audit.GetTangibleAuditLog)

		// 用户管理&user
		auth.GET("/users", user.GetSysUserList)
		auth.POST("/users", user.CreateUser)
		auth.GET("/userinfo", user.GetInfo)
		auth.GET("/user/profile", user.GetInfo)
		auth.PUT("/users/:userID", user.UpdateUser)
		auth.PUT("/user/profile", user.UpdateUserProfile)
		auth.POST("/user/upload", user.ImportUserData)
		auth.GET("/users/:userID", user.GetSysUser)
		auth.DELETE("/users/:userID", user.DeleteSysUser)

		// 角色管理&role
		auth.GET("/roles", role.GetSysRoleList)
		auth.GET("/allroles", role.GetAllRoles)
		auth.POST("/roles", role.CreateRole)
		auth.GET("/roles/:roleID", role.GetRole)
		auth.PUT("/roles/:roleID", role.UpdateRole)
		auth.DELETE("/roles/:roleID", role.DeleteRole)

		// 菜单管理&menu
		auth.GET("/menuTreeselect", sysadmin.GetMenuTreeelect)
		auth.GET("/menuList", sysadmin.GetMenuList)
		auth.POST("/menus", sysadmin.InsertMenu)
		auth.PUT("/menus/:id", sysadmin.UpdateMenu)
		auth.GET("/menus/:id", sysadmin.GetMenu)
		auth.DELETE("/menus/:id", sysadmin.DeleteMenu)
		auth.GET("/roleMenuTreeselect/:roleID", sysadmin.GetMenuTreeRoleselect)
		auth.GET("/role/menus", sysadmin.GetRoleMenu)
	}
	log.Println("router load success")
	return r
}
