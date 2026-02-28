package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Refresh 处理 Token 刷新 POST api/user/account/refresh_token/
func (h *Handler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": "refresh token 不存在",
		})
		return
	}

	newAccess, newRefresh, err := h.authSvc.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"result": "refresh token 过期了",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", newRefresh, 86400*7, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"access": newAccess,
	})
}
