package mw

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/aikwen/aifriend-go/pkg/auth"
)

// JWTAuthMiddleware 拦截并校验请求头中的 JWT 令牌
func JWTAuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"result": "未提供认证信息"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"result": "认证头格式错误，请使用 Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := auth.ParseAccessToken(tokenString, []byte(secret))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"result": "Token 无效或已过期"})
			c.Abort()
			return
		}

		userID64, err := strconv.ParseUint(claims.Subject, 10, 32)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"result": "Token 载荷异常"})
			c.Abort()
			return
		}

		c.Set("userID", uint(userID64))
		c.Set("username", claims.Name)

		c.Next()
	}
}