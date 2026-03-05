package mw

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/aikwen/aifriend-go/config"
)

// ServiceReadyCheck 全局服务就绪检查中间件
func ServiceReadyCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果配置中关闭了服务，则拦截所有请求
		if !config.GlobalConfig.Server.Enable {
			c.JSON(http.StatusOK, gin.H{
				"result": "服务器维护中，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}