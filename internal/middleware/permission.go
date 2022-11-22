package middleware

import (
	"net/http"
	"strings"

	mycasbin "github.com/colynn/pontus/internal/middleware/casbin"
	"github.com/colynn/pontus/internal/middleware/jwtauth"
	"github.com/colynn/pontus/internal/pkg/customerror"

	"github.com/gin-gonic/gin"
	log "unknwon.dev/clog/v2"
)

// AuthCheckRole ..
func AuthCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, _ := c.Get("JWT_PAYLOAD")
		v := data.(jwtauth.MapClaims)
		e, err := mycasbin.NewCasbin()
		customerror.HasError(err, "", 500)
		//检查权限
		URLPath := c.Request.URL.Path
		URLMethod := c.Request.Method

		if URLMethod == "GET" {
			return
		}

		// !strings.Contains(URLPath, "api/v1/roles") || !strings.Contains(URLPath, "api/v1/menus") || !strings.Contains(URLPath, "api/v1/dept") || !strings.Contains(URLPath, "api/v1/users")

		// current if roleKey is admin, skip permission verify
		// TODO: sql audit api/v2 skip api permission verify
		if v["rolekey"] == "admin" || strings.Contains(URLPath, "api/v2/") {
			return
		}
		res, err := e.Enforce(v["rolekey"], URLPath, URLMethod)
		log.Trace("role key: %s, path: %s, method: %s, res: %v", v["rolekey"], c.Request.URL.Path, c.Request.Method, res)

		customerror.HasError(err, "", 500)

		if res {
			c.Next()
		} else {
			c.JSON(http.StatusOK, gin.H{
				"code": 403,
				"msg":  "对不起，您没有访问权限，请联系管理员",
			})
			c.Abort()
			return
		}
	}
}
