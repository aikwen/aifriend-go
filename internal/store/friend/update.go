package friend

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/errs"
	"github.com/aikwen/aifriend-go/internal/store/models"
)

func (f *friendStore) UpdateMemory(ctx context.Context, userID uint, friendID uint, memory string) error {
	subQuery := f.db.WithContext(ctx).
		Model(&models.Friend{}).
		Select("friends.id").
		Joins("INNER JOIN characters ON characters.id = friends.character_id AND characters.deleted_at IS NULL").
		Where("friends.id = ? AND friends.user_id = ? AND friends.deleted_at IS NULL", friendID, userID)

	result := f.db.WithContext(ctx).
		Model(&models.Friend{}).
		Where("id IN (?)", subQuery).
		Update("memory", memory)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errs.ErrFriendNotFound
	}

	return nil
}