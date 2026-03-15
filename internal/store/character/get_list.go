package character

import (
	"context"
	"strings"

	"github.com/aikwen/aifriend-go/internal/store/models"
	"gorm.io/gorm"
)

// GetListByAuthorID 根据用户id和 offset 获取 limit 条数据
func (s *characterStore) GetListByAuthorID(ctx context.Context, authorID uint, offset int, limit int) ([]*models.Character, error) {
	var characters []*models.Character
	err := s.db.WithContext(ctx).
		Where("author_id = ?", authorID).
		Preload("Author").
		Order("id desc").
		Offset(offset).
		Limit(limit).
		Find(&characters).Error
	return characters, err
}

// getListbySearchQuery 根据关键字搜索虚拟角色
func (s *characterStore) GetListBySearchQuery(ctx context.Context, offset int, limit int, searchQuery string) ([]*models.Character, error) {
	var characters []*models.Character

	searchQuery = strings.TrimSpace(searchQuery)

	query := s.db.WithContext(ctx).Preload("Author")

	if searchQuery != "" {
		query = query.Where("MATCH (name, profile) AGAINST (? IN NATURAL LANGUAGE MODE)", searchQuery).
                  Order(gorm.Expr("MATCH (name, profile) AGAINST (? IN NATURAL LANGUAGE MODE) DESC", searchQuery)).
				  Order("id DESC")
	} else {
		query = query.Order("id DESC")
	}

	err := query.
		Offset(offset).
		Limit(limit).
		Find(&characters).
		Error

	return characters, err
}