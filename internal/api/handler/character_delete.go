package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DeleteCharacter 处理删除角色的 POST 请求 /api/create/character/remove/
func (h *Handler) DeleteCharacter(c *gin.Context) {
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

	var req struct {
		CharacterID uint `json:"character_id" form:"character_id"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
		return
	}

	if err := h.charSvc.DeleteCharacter(c.Request.Context(), req.CharacterID, userID); err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
