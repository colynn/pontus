package accountinfo

import (
	"strconv"

	jwt "github.com/colynn/pontus/internal/middleware/jwtauth"

	log "unknwon.dev/clog/v2"

	"github.com/gin-gonic/gin"
)

// ExtractClaims ..
func ExtractClaims(c *gin.Context) jwt.MapClaims {
	claims, exists := c.Get("JWT_PAYLOAD")
	if !exists {
		return make(jwt.MapClaims)
	}

	return claims.(jwt.MapClaims)
}

// GetUserID ..
func GetUserID(c *gin.Context) int {
	data := ExtractClaims(c)
	if data["identity"] != nil {
		return int((data["identity"]).(float64))
	}
	log.Warn("******** url.path: " + c.Request.URL.Path + "  request.method: " + c.Request.Method + "  说明：缺少identity")
	return 0
}

// GetUserIDStr ..
func GetUserIDStr(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["identity"] != nil {
		return strconv.FormatInt(int64((data["identity"]).(float64)), 10)
	}
	log.Warn("******** url.path: " + c.Request.URL.Path + "  request.method: " + c.Request.Method + "  缺少identity")
	return ""
}

// GetUserName ..
func GetUserName(c *gin.Context) string {
	data := ExtractClaims(c)
	if data["username"] != nil {
		return (data["username"]).(string)
	}
	log.Warn("******** url.path: " + c.Request.URL.Path + "  request.method: " + c.Request.Method + "  缺少username")
	return ""
}
