package handler

import (
	"net/http"
	"strings"

	"github.com/aikwen/aifriend-go/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(as service.AuthService) *AuthHandler {
	return &AuthHandler{
		authSvc: as,
	}
}

// Register 处理用户注册 POST api/user/account/register/
func (ah *AuthHandler) Register(c *gin.Context) {
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

	user, access, refresh, err := ah.authSvc.Register(c.Request.Context(), username, password)
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

// Login 处理用户登录 POST api/user/account/login/
func (h *AuthHandler) Login(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, gin.H{"result": err.Error()})
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
func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"result": "success",
	})
}


// Refresh 处理 Token 刷新 POST api/user/account/refresh_token/
func (h *AuthHandler) Refresh(c *gin.Context) {
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