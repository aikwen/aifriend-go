package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FriendGetOrCreate 处理获取或创建好友关系的请求
func (h *Handler) GetOrCreateFriend(c *gin.Context) {
	var req struct {
		CharacterID uint `json:"character_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "参数错误"})
		return
	}

	userID := c.MustGet("userID").(uint)

	friendDTO, err := h.friendSvc.GetOrCreate(userID, req.CharacterID)
	if err != nil {
		log.Printf("[FriendGetOrCreate] failed: %v\n", err)
		c.JSON(http.StatusOK, gin.H{"result": "系统异常，请稍后再试"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"friend": friendDTO,
	})
}