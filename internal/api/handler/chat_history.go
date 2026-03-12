package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)


func (h *Handler) GetChatHistory(c *gin.Context) {
	friendIDStr := c.Query("friend_id")
	lastMessageIDStr := c.Query("last_message_id")

	friendID, err := strconv.ParseUint(friendIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "参数错误"})
		return
	}

	var lastMessageID uint64
	if lastMessageIDStr != "" {
		lastMessageID, _ = strconv.ParseUint(lastMessageIDStr, 10, 64)
	}

	userID := c.MustGet("userID").(uint)

	messages, err := h.chatSvc.GetHistory(c.Request.Context(), uint(friendID), uint(lastMessageID), userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "服务错误"})
		return
	}

	// 组装消息
	respMessages := []gin.H{}
	for _, m := range messages {
		respMessages = append(respMessages, gin.H{
			"id":           m.ID,
			"user_message": m.UserMessage,
			"output":       m.Output,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"result":   "success",
		"messages": respMessages,
	})
}