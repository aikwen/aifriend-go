package character

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
)


// getListByAuthorID 根据用户id和 offset 获取 limit 条数据
func (s *characterStore) getListByAuthorID(ctx context.Context, authorID uint, offset int, limit int) ([]*models.Character, error) {
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
func (s *characterStore) getListBySearchQuery(ctx context.Context, offset int, limit int, searchQuery string) ([]*models.Character, error) {
	var characters []*models.Character

	query := s.db.WithContext(ctx).Preload("Author")

	if searchQuery != "" {
		likePattern := "%" + searchQuery + "%"
		query = query.Where("name LIKE ? OR profile LIKE ?", likePattern, likePattern)
	}

	err := query.
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&characters).
		Error

	return characters, err
}