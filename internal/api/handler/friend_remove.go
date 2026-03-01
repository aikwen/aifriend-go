package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FriendRemove 处理删除好友的请求
func (h *Handler) RemoveFriend(c *gin.Context) {
	var req struct {
		FriendID uint `json:"friend_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "参数错误"})
		return
	}

	userID := c.MustGet("userID").(uint)

	err := h.friendSvc.Remove(userID, req.FriendID)
	if err != nil {
		log.Printf("[FriendRemove] failed: %v\n", err)
		c.JSON(http.StatusOK, gin.H{"result": "系统异常，请稍后重试"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "success",
	})
}