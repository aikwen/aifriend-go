package character

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
)


// Update 更新角色信息
func (s *characterStore) update(ctx context.Context, character *models.Character) error {
	return s.db.WithContext(ctx).Save(character).Error
}