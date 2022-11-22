package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/colynn/pontus/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//实例化
var logger = logrus.New()

// LoggerToFile 日志记录到文件
func LoggerToFile() gin.HandlerFunc {

	c := config.GetConfig()
	//写入文件
	// TODO: confirm log.auditlog exists
	src, err := os.OpenFile(c.GetString("log.auditlog"), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println("err", err)
	}

	//设置输出
	logger.Out = src

	//设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	//设置日志格式
	logger.SetFormatter(&logrus.TextFormatter{})

	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqURI := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logger.Infof(" %3d  %13v  %15s  %s  %s \r\n",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqURI,
		)
		fmt.Println(
			startTime.Format("2006-01-02 15:04:05.9999"),
			"[INFO]",
			reqMethod,
			reqURI,
			statusCode,
			latencyTime,
			reqURI,
			clientIP,
		)

		// if c.Request.Method != "GET" && c.Request.Method != "OPTIONS" {
		// 	menu := models.Menu{}
		// 	menu.Path = reqUri
		// 	menu.Action = reqMethod
		// 	menuList, _ := menu.Get()
		// 	sysOperLog := user.SysOperLog{}
		// 	sysOperLog.OperIp = clientIP
		// 	sysOperLog.OperLocation = utils.GetLocation(clientIP)
		// 	sysOperLog.Status = utils.IntToString(statusCode)
		// 	sysOperLog.OperName = utils.GetUserName(c)
		// 	sysOperLog.RequestMethod = c.Request.Method
		// 	sysOperLog.OperUrl = reqUri
		// 	if reqUri == "/login" {
		// 		sysOperLog.BusinessType = "10"
		// 		sysOperLog.Title = "用户登录"
		// 		sysOperLog.OperName = "-"
		// 	} else if strings.Contains(reqUri, "/api/v1/logout") {
		// 		sysOperLog.BusinessType = "11"
		// 	} else if strings.Contains(reqUri, "/api/v1/getCaptcha") {
		// 		sysOperLog.BusinessType = "12"
		// 		sysOperLog.Title = "验证码"
		// 	} else {
		// 		if reqMethod == "POST" {
		// 			sysOperLog.BusinessType = "1"
		// 		} else if reqMethod == "PUT" {
		// 			sysOperLog.BusinessType = "2"
		// 		} else if reqMethod == "DELETE" {
		// 			sysOperLog.BusinessType = "3"
		// 		}
		// 	}
		// 	sysOperLog.Method = reqMethod
		// 	if len(menuList) > 0 {
		// 		sysOperLog.Title = menuList[0].Title
		// 	}
		// 	b, _ := c.Get("body")
		// 	sysOperLog.OperParam, _ = utils.StructToJsonStr(b)
		// 	sysOperLog.CreateBy = utils.GetUserName(c)
		// 	sysOperLog.OperTime = utils.GetCurrntTime()
		// 	sysOperLog.LatencyTime = (latencyTime).String()
		// 	sysOperLog.UserAgent = c.Request.UserAgent()
		// 	if c.Err() == nil {
		// 		sysOperLog.Status = "0"
		// 	} else {
		// 		sysOperLog.Status = "1"
		// 	}
		// 	_, _ = sysOperLog.Create()
		// }
	}
}
