package friend

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
)


func (s *friendStore) getOrCreate(ctx context.Context, userID uint, characterID uint) (*models.Friend, error) {
	var friend models.Friend

	// 根据 Where 里的条件去查，查不到就用 Where 里的条件做初始化并插入
	err := s.db.WithContext(ctx).Where(models.Friend{UserID: userID, CharacterID: characterID}).
		FirstOrCreate(&friend).Error
	if err != nil {
		return nil, err
	}

	err = s.db.WithContext(ctx).Preload("Character").
		Preload("Character.Author").
		First(&friend, friend.ID).Error

	return &friend, err
}