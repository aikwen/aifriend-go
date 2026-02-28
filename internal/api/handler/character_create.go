package handler

import (
	"fmt"
	"net/http"

	"github.com/aikwen/aifriend-go/internal/character"
	"github.com/gin-gonic/gin"
)

// CreateCharacter 处理创建角色的 POST 请求 /api/create/character/create/
func (h *Handler) CreateCharacter(c *gin.Context) {
	// 从 context 获取 userID
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

	// 提取表单参数
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

	// 提取文件参数
	photoHeader, err1 := c.FormFile("photo")
	bgImageHeader, err2 := c.FormFile("background_image")

	if err1 != nil || photoHeader == nil {
		c.JSON(http.StatusOK, gin.H{"result": "头像不能为空"})
		return
	}
	if err2 != nil || bgImageHeader == nil {
		c.JSON(http.StatusOK, gin.H{"result": "聊天背景不能为空"})
		return
	}

	// 保存文件
	userIDStr := fmt.Sprintf("%d", userID)
	photoPath, err := h.storage.Save(photoHeader, "character/photos", userIDStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "头像保存失败"})
		return
	}

	bgImagePath, err := h.storage.Save(bgImageHeader, "character/background_images", userIDStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "背景图保存失败"})
		return
	}

	// 调用 Service
	param := &character.CreateCharacterParam{
		AuthorID:    userID,
		Name:        name,
		Profile:     profile,
		PhotoPath:   photoPath,
		BgImagePath: bgImagePath,
	}

	if err := h.charSvc.CreateCharacter(c.Request.Context(), param); err != nil {
		c.JSON(http.StatusOK, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}
