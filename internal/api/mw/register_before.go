package mw

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/aikwen/aifriend-go/config"
)

// RegisterBeforeCheck 注册前置校验中间件
func RegisterBeforeCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查全局配置开关
		if !config.GlobalConfig.Server.AllowRegister {
			c.JSON(http.StatusOK, gin.H{
				"result": "当前系统已关闭注册功能",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}