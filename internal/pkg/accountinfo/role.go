package accountinfo

import (
	"github.com/gin-gonic/gin"
	log "unknwon.dev/clog/v2"
)

// GetRoleID ..
func GetRoleID(c *gin.Context) int {
	data := ExtractClaims(c)
	if data["roleid"] != nil {
		return int((data["roleid"]).(float64))
	}
	log.Warn("*************** url.path" + c.Request.URL.Path + "  request.method: " + c.Request.Method + "  说明：缺少roleid")
	return 0
}

// GetRoleName ..
func GetRoleName(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["rolekey"] != nil {
		return (data["rolekey"]).(string)
	}
	log.Warn("**************** url.path: " + c.Request.URL.Path + "  request.method: " + c.Request.Method + "  缺少rolekey")
	return ""
}
