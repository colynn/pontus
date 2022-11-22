package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/colynn/pontus/config"
	"github.com/colynn/pontus/internal/account/role"
	"github.com/colynn/pontus/internal/account/user"
	jwt "github.com/colynn/pontus/internal/middleware/jwtauth"
	"github.com/colynn/pontus/tools/currenttime"
	"github.com/colynn/pontus/tools/locationip"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/mssola/user_agent"
	log "unknwon.dev/clog/v2"
)

var store = base64Captcha.DefaultMemStore

// AuthInit ..
func AuthInit() (*jwt.GinJWTMiddleware, error) {
	c := config.GetConfig()
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "jwt auth",
		Key:             []byte(c.GetString("jwt.secret")),
		Timeout:         time.Duration(c.GetInt64("jwt.timeout")) * time.Hour,
		MaxRefresh:      time.Duration(c.GetInt64("jwt.timeout")) * time.Hour,
		PayloadFunc:     PayloadFunc,
		IdentityHandler: IdentityHandler,
		Authenticator:   Authenticator,
		Unauthorized:    Unauthorized,
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})
}

// PayloadFunc ..
func PayloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		u, _ := v["user"].(*user.SysUser)
		r, _ := v["role"].(*role.SysRole)
		return jwt.MapClaims{
			jwt.IdentityKey: u.SysUserId.ID,
			jwt.RoleKey:     r.RoleKey,
			jwt.UserName:    u.Username,
		}
	}
	return jwt.MapClaims{}
}

// IdentityHandler ..
func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return map[string]interface{}{
		"IdentityKey": claims["identity"],
		"UserName":    claims["username"],
		"RoleKey":     claims["rolekey"],
	}
}

// Authenticator ..
func Authenticator(c *gin.Context) (interface{}, error) {
	var loginVals user.Login
	var loginlog user.LoginLog
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
	loginlog.Msg = "登录成功"
	loginlog.Platform = ua.Platform()

	if err := c.ShouldBind(&loginVals); err != nil {
		loginlog.Status = "1"
		loginlog.Msg = "数据解析失败"
		loginlog.Username = loginVals.Username
		loginlog.Create()
		if strings.Contains(err.Error(), "Login.UUID' Error:Field validation for 'UUID' failed on the 'required' tag") {
			return nil, jwt.ErrMissingUUID
		}
		return nil, jwt.ErrMissingLoginValues
	}
	loginlog.Username = loginVals.Username
	// Comment 验证码 tmp
	// if !store.Verify(loginVals.UUID, loginVals.Code, true) {
	// 	loginlog.Status = "1"
	// 	loginlog.Msg = "验证码错误"
	// 	loginlog.Create()
	// 	return nil, jwt.ErrInvalidVerificationode
	// }

	userSvc := user.NewService()
	req := &user.Login{
		Username: loginVals.Username,
		Password: loginVals.Password,
	}
	userInfo, err := userSvc.Authenticate(req)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	log.Trace("ldap auth success: %v", userInfo.Name)
	exist, userData, e := userSvc.Exists(&loginVals)
	if e == nil {
		roleModel := role.SysRole{}
		if !exist {
			// TODO: default role get from config
			roleModel.RoleKey = "admin"
			roleItem, err := roleModel.Get()
			if err != nil {
				log.Error("when create user get role by role key error: %s", err.Error())
				return nil, jwt.ErrFailedAuthentication
			}
			user := &user.SysUser{
				SysUserB: user.SysUserB{
					Email:    userInfo.Email,
					RealName: userInfo.Name,
					RoleID:   roleItem.RoleID,
				},
				LoginM: user.LoginM{
					Username: loginVals.Username,
				},
			}
			id, err := userSvc.UserDao.Insert(user)
			if err != nil {
				log.Error("user insert error: %s", err.Error())
			}
			log.Trace("user id: %d", id)
			userData, err = userSvc.UserDao.Get(id)
			if err != nil {
				log.Error("user get error: %s", err.Error())
			}
		}
		roleModel.RoleID = userData.RoleID
		roleData, err := roleModel.Get()
		if err != nil {
			log.Error("get role occur error: %s", err.Error())
			// assign empty value
			roleData = roleModel
		}
		loginlog.Create()
		log.Trace("logined user: %+v", userData)
		return map[string]interface{}{"user": userData, "role": &roleData}, nil
	}

	loginlog.Status = "1"
	loginlog.Msg = "登录失败"
	loginlog.Create()
	log.Error("login faile, error: %s", e.Error())

	return nil, jwt.ErrForbidden
}

// Unauthorized ..
func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  message,
	})
}
