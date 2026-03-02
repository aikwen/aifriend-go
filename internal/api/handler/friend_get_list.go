package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// FriendGetList 处理获取好友列表的请求
func (h *Handler) GetFriendList(c *gin.Context) {

	itemsCountStr := c.DefaultQuery("items_count", "0")
	offset, err := strconv.Atoi(itemsCountStr)
	if err != nil || offset < 0 {
		offset = 0 // 如果解析失败，默认从头开始
	}

	limit := 20

	userID := c.MustGet("userID").(uint)

	friendsList, err := h.friendSvc.GetList(c.Request.Context(), userID, offset, limit)
	if err != nil {
		log.Printf("[FriendGetList] failed: %v\n", err)
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":  "success",
		"friends": friendsList,
	})
}