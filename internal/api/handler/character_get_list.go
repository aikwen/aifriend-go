package handler

import (
	"net/http"
	"path"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCharacterList 获取角色列表 get /api/create/character/get_list/
func (h *Handler) GetCharacterList(c *gin.Context) {
	itemsCountStr := c.Query("items_count")
	userIDStr := c.Query("user_id")

	itemsCount, err := strconv.Atoi(itemsCountStr)
	if err != nil {
		itemsCount = 0
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
		return
	}

	user, err := h.userSvc.GetByID(c.Request.Context(), uint(userID))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
		return
	}

	characters, err := h.charSvc.GetCharacterList(c.Request.Context(), uint(userID), itemsCount)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"result": "系统异常"})
		return
	}

	var charDataList []gin.H
	for _, char := range characters {
		charDataList = append(charDataList, gin.H{
			"id":               char.ID,
			"name":             char.Name,
			"profile":          char.Profile,
			"photo":            path.Join("/media/", char.Photo),
			"background_image": path.Join("/media/", char.BackgroundImage),
			"author": gin.H{
				"user_id":  char.Author.ID,
				"username": char.Author.Username,
				"photo":    path.Join("/media/", char.Author.Photo),
			},
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"result": "success",
		"user_profile": gin.H{
			"user_id":  user.ID,
			"username": user.Username,
			"profile":  user.Profile,
			"photo":    path.Join("/media/", user.Photo),
		},
		"characters": charDataList,
	})
}