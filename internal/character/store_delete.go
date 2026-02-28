package character


import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
)


// delete 删除角色，
func (s *characterStore) delete(ctx context.Context, id uint, authorID uint) error {
	return s.db.WithContext(ctx).
		Where("id = ? AND author_id = ?", id, authorID).
		Delete(&models.Character{}).Error
}