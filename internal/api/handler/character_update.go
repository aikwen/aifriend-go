package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/aikwen/aifriend-go/internal/character"
	"github.com/gin-gonic/gin"
)

// UpdateCharacter 处理更新角色的 POST 请求 /api/create/character/update/
func (h *Handler) UpdateCharacter(c *gin.Context) {
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

	charIDStr := c.PostForm("character_id")
	charID, err := strconv.ParseUint(charIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
		return
	}

	name := c.PostForm("name")
	profile := c.PostForm("profile")

	if name == "" {
		c.JSON(http.StatusOK, gin.H{"result": "名字不能为空"})
		return
	}
	if profile == "" {
		c.JSON(http.StatusOK, gin.H{"result": "角色介绍不能为空"})
		return
	}

	photoPath := ""
	userIDStr := fmt.Sprintf("%d", userID)
	if photoHeader, err := c.FormFile("photo"); err == nil && photoHeader != nil {
		photoPath, err = h.storage.Save(photoHeader, "character/photos", userIDStr)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"result": "头像保存失败"})
			return
		}
	}

	bgImagePath := ""
	if bgImageHeader, err := c.FormFile("background_image"); err == nil && bgImageHeader != nil {
		bgImagePath, err = h.storage.Save(bgImageHeader, "character/background_images", userIDStr)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{"result": "背景图保存失败"})
			return
		}
	}

	param := &character.UpdateCharacterParam{
		ID:          uint(charID),
		AuthorID:    userID,
		Name:        name,
		Profile:     profile,
		PhotoPath:   photoPath,
		BgImagePath: bgImagePath,
	}

	if err := h.charSvc.UpdateCharacter(c.Request.Context(), param); err != nil {
		c.JSON(http.StatusOK, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
