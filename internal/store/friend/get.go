package friend

import (
	"context"
	"errors"

	"github.com/aikwen/aifriend-go/internal/errs"
	"github.com/aikwen/aifriend-go/internal/store/models"
	"gorm.io/gorm"
)

// 查询用户可访问的好友及其角色信息
func (f *friendStore) GetByIDAndUserID(ctx context.Context, userID uint, friendID uint) (*models.Friend, error) {
	var friend models.Friend

	err := f.db.WithContext(ctx).
		Model(&models.Friend{}).
		Preload("Character").
		Joins("INNER JOIN characters ON characters.id = friends.character_id AND characters.deleted_at IS NULL").
		Where("friends.id = ? AND friends.user_id = ? AND friends.deleted_at IS NULL", friendID, userID).
		First(&friend).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.ErrFriendNotFound
		}
		return nil, err
	}

	return &friend, nil
}