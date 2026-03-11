package character

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
)

// Delete 删除角色，
func (s *characterStore) Delete(ctx context.Context, id uint, authorID uint) error {
	return s.db.WithContext(ctx).
		Where("id = ? AND author_id = ?", id, authorID).
		Delete(&models.Character{}).Error
}
