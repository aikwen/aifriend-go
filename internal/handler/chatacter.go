package handler

import (
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/aikwen/aifriend-go/internal/service"
	"github.com/aikwen/aifriend-go/pkg/storage"
	"github.com/gin-gonic/gin"
)

type CharacterHandler struct {
	CharacterSrc service.CharacterService
	storage storage.Storage
}

func NewCharacterHandler(cs service.CharacterService, st storage.Storage) *CharacterHandler {
	return &CharacterHandler{
		CharacterSrc: cs,
		storage: st,
	}
}


// CreateCharacter 处理创建角色的 POST 请求 /api/create/character/create/
func (h *CharacterHandler) CreateCharacter(c *gin.Context) {
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
	param := &service.CreateCharacterParam{
		AuthorID:    userID,
		Name:        name,
		Profile:     profile,
		PhotoPath:   photoPath,
		BgImagePath: bgImagePath,
	}

	if err := h.CharacterSrc.CreateCharacter(c.Request.Context(), param); err != nil {
		c.JSON(http.StatusOK, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

// GetCharacter 处理获取单个角色的 GET 请求 /api/create/character/get_single/
func (h *CharacterHandler) GetCharacter(c *gin.Context) {
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

	char, err := h.CharacterSrc.GetCharacter(c.Request.Context(), uint(charID), userID)
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

// UpdateCharacter 处理更新角色的 POST 请求 /api/create/character/update/
func (h *CharacterHandler) UpdateCharacter(c *gin.Context) {
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

	param := &service.UpdateCharacterParam{
		ID:          uint(charID),
		AuthorID:    userID,
		Name:        name,
		Profile:     profile,
		PhotoPath:   photoPath,
		BgImagePath: bgImagePath,
	}

	if err := h.CharacterSrc.UpdateCharacter(c.Request.Context(), param); err != nil {
		c.JSON(http.StatusOK, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

// DeleteCharacter 处理删除角色的 POST 请求 /api/create/character/remove/
func (h *CharacterHandler) DeleteCharacter(c *gin.Context) {
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

	if err := h.CharacterSrc.DeleteCharacter(c.Request.Context(), req.CharacterID, userID); err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "success"})
}

