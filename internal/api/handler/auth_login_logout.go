package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Login 处理用户登录 POST api/user/account/login/
func (h *Handler) Login(c *gin.Context) {
	req := struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": "参数错误"})
		return
	}

	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)

	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"result": "用户名和密码不能为空",
		})
		return
	}

	user, accessToken, refreshToken, err := h.authSvc.Login(c.Request.Context(), username, password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode) // samesite='Lax'
	c.SetCookie("refresh_token", refreshToken, 86400*7, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"result":   "success",
		"access":   accessToken, // 注意这里是 access
		"user_id":  user.ID,
		"username": user.Username,
		"photo":    "/media/" + user.Photo,
		"profile":  user.Profile,
	})
}

// Logout 处理用户登出 POST api/user/account/logout/
func (h *Handler) Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"result": "success",
	})
}
