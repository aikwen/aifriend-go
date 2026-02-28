package character


import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
)


// create 创建角色
func (s *characterStore) create(ctx context.Context, character *models.Character) error {
	return s.db.WithContext(ctx).Create(character).Error
}
