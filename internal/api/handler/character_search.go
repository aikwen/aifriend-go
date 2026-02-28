package handler

import (
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// SearchCharacters 根据搜索内容搜索角色 get /api/homepage/index/
func (h *Handler) SearchCharacters(c *gin.Context) {

	itemsCountStr := c.Query("items_count")
	searchQuery := strings.TrimSpace(c.Query("search_query"))


	offset, err := strconv.Atoi(itemsCountStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	limit := 20


	characters, err := h.charSvc.SearchCharacters(c.Request.Context(), offset, limit, searchQuery)
	if err != nil {
		log.Printf("[SearchCharacters] 业务查询失败, offset: %d, search: %s, err: %v", offset, searchQuery, err)
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

	if charDataList == nil {
		charDataList = make([]gin.H, 0)
	}

	c.JSON(http.StatusOK, gin.H{
		"result":     "success",
		"characters": charDataList,
	})
}