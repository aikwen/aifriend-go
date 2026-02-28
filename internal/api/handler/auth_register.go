package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Register 处理用户注册 POST api/user/account/register/
func (h *Handler) Register(c *gin.Context) {
	req := struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}{}

	// 自动解析 JSON 并校验
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": "参数错误", "error": err.Error()})
		return
	}

	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)

	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{"result": "用户名或者密码不能为空"})
		return
	}

	user, access, refresh, err := h.authSvc.Register(c.Request.Context(), username, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": "注册失败", "message": err.Error()})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode) // 设置 SameSite 为 Lax
	c.SetCookie("refresh_token", refresh, 86400*7, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"result":   "success",
		"access":   access,
		"user_id":  user.ID,
		"username": user.Username,
		"photo":    "/media/" + user.Photo,
		"profile":  user.Profile,
	})
}
