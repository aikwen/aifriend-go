package handler

import (
	"fmt"
	"net/http"
	"path"

	"github.com/aikwen/aifriend-go/internal/service"
	"github.com/aikwen/aifriend-go/pkg/storage"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc service.UserService
	storage storage.Storage
}

func NewUserHandler(us service.UserService, st storage.Storage) *UserHandler {
	return &UserHandler{
		userSvc: us,
		storage: st,
	}
}

// GetUserInfo 处理获取用户信息的 GET /api/user/account/get_user_info/
// protect
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	// 从context 获取 userID
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
		return
	}

	userID, ok := userIDAny.(uint)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
		return
	}

	user, err := h.userSvc.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":   "success",
		"user_id":  user.ID,
		"username": user.Username,
		"photo":    path.Join("/media/", user.Photo),
		"profile":  user.Profile,
	})
}

// UpdateUserInfo 更新用户信息 post /api/user/profile/update/
// protect
func (h *UserHandler) UpdateUserInfo(c *gin.Context) {
	userIDAny, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
		return
	}

	userID, ok := userIDAny.(uint)
	if !ok {
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
		return
	}

	// 提取参数
	username := c.PostForm("username")
	profile := c.PostForm("profile")

	if username == "" {
		c.JSON(http.StatusOK, gin.H{"result": "用户名不能为空"})
		return
	}
	if profile == "" {
		c.JSON(http.StatusOK, gin.H{"result": "简介不能为空"})
		return
	}

	fileHeader, err := c.FormFile("photo")
	photoPath := ""
	if err == nil && fileHeader != nil {
		userIDStr := fmt.Sprintf("%d", userID)
		photoPath, err = h.storage.Save(fileHeader, "user/photos", userIDStr)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"result": "头像保存失败，请稍后重试"})
			return
		}
	}

	updatedUser, err := h.userSvc.UpdateUserInfo(c.Request.Context(), userID, username, profile, photoPath)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":   "success",
		"user_id":  updatedUser.ID,
		"username": updatedUser.Username,
		"profile":  updatedUser.Profile,
		"photo":    path.Join("/media/", updatedUser.Photo),
	})
}
