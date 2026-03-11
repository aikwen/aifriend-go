package character

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
)

// Create 创建角色
func (s *characterStore) Create(ctx context.Context, character *models.Character) error {
	return s.db.WithContext(ctx).Create(character).Error
}
