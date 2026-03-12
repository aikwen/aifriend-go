package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/aikwen/aifriend-go/internal/chat"
	"github.com/gin-gonic/gin"
)

// ai 聊天
func (h *Handler) Chat(c *gin.Context) {
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

	type chatReq struct {
		FriendID uint   `json:"friend_id"`
		Message  string `json:"message"`
	}

	var req chatReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"detail": "请求参数错误"})
		return
	}

	req.Message = strings.TrimSpace(req.Message)
	if req.FriendID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "好友参数错误"})
		return
	}
	if req.Message == "" {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "消息不能为空"})
		return
	}

	// 获取事件流
	events, err := h.chatSvc.Chat(c.Request.Context(), userID, req.FriendID, req.Message)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": err.Error()})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	c.Stream(func(w io.Writer) bool {
		event, ok := <-events
		if !ok {
			return false
		}

		switch event.Type {
		case chat.EventDelta:
			if event.Text == "" {
				return true
			}

			data, err := json.Marshal(gin.H{
				"content": event.Text,
			})
			if err != nil {
				return false
			}

			_, _ = w.Write([]byte("data: "))
			_, _ = w.Write(data)
			_, _ = w.Write([]byte("\n\n"))
			return true

		case chat.EventDone:
			_, _ = w.Write([]byte("data: [DONE]\n\n"))
			return false

		case chat.EventError:
			return false

		default:
			return true
		}
	})
}
