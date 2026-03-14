package friend

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
)

func (f *friendStore) Exists(ctx context.Context, userID uint, friendID uint) (bool, error) {
	var count int64

	err := f.db.WithContext(ctx).
		Model(&models.Friend{}).
		Joins("INNER JOIN characters ON characters.id = friends.character_id AND characters.deleted_at IS NULL").
		Where("friends.id = ? AND friends.user_id = ? AND friends.deleted_at IS NULL", friendID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}