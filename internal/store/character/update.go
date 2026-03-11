package character

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
)

// Update 更新角色信息
func (s *characterStore) Update(ctx context.Context, character *models.Character) error {
	return s.db.WithContext(ctx).Save(character).Error
}
