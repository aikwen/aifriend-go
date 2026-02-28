package handler

import (
	"net/http"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCharacter 处理获取单个角色的 GET 请求 /api/create/character/get_single/
func (h *Handler) GetCharacter(c *gin.Context) {
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

	charIDStr := c.Query("character_id")
	charID, err := strconv.ParseUint(charIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
		return
	}

	char, err := h.charSvc.GetCharacter(c.Request.Context(), uint(charID), userID)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"character": gin.H{
			"id":               char.ID,
			"name":             char.Name,
			"profile":          char.Profile,
			"photo":            path.Join("/media/", char.Photo),
			"background_image": path.Join("/media/", char.BackgroundImage),
		},
	})
}
