package character

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
)

// getByIDAndAuthor 获取单个角色
func (s *characterStore) getByIDAndAuthor(ctx context.Context, id uint, authorID uint) (*models.Character, error) {
	var c models.Character
	// 获取对应的authorId下的character id
	err := s.db.WithContext(ctx).
		Where("id = ? AND author_id = ?", id, authorID).
		First(&c).Error

	if err != nil {
		return nil, err
	}
	return &c, nil
}