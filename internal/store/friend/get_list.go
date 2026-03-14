package friend

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
)

func (s *friendStore) GetList(ctx context.Context, userID uint, offset int, limit int) ([]models.Friend, error) {
	var friends []models.Friend

	err := s.db.WithContext(ctx).
		Model(&models.Friend{}).
		Joins("INNER JOIN characters ON characters.id = friends.character_id AND characters.deleted_at IS NULL").
		Where("friends.user_id = ? AND friends.deleted_at IS NULL", userID).
		Order("friends.updated_at DESC").
		Offset(offset).
		Limit(limit).
		Preload("Character.Author").
		Find(&friends).Error

	return friends, err
}