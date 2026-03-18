package friend

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/errs"
	"github.com/aikwen/aifriend-go/internal/store/models"
)

func (f *friendStore) UpdateMemoryWithVersion(ctx context.Context, userID uint, friendID uint, oldVersion uint, memory string) error {
	result := f.db.WithContext(ctx).
		Model(&models.Friend{}).
		Joins("INNER JOIN characters ON characters.id = friends.character_id AND characters.deleted_at IS NULL").
		Where("friends.id = ? AND friends.user_id = ? AND friends.deleted_at IS NULL AND friends.version = ?", friendID, userID, oldVersion).
		Updates(map[string]any{
			"memory":  memory,
			"version": oldVersion + 1,
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errs.ErrFriendVersionConflict
	}

	return nil
}