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