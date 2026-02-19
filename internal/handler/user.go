package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/aikwen/aifriend-go/internal/service"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(us service.UserService) *UserHandler {
	return &UserHandler{
		userSvc: us,
	}
}

// GetUserInfo 处理获取用户信息的 GET /api/user/account/get_user_info/
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
		"photo":    "/media/" + user.Photo,
		"profile":  user.Profile,
	})
}