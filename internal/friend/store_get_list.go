package friend

import "github.com/aikwen/aifriend-go/internal/models"


func (s *friendStore) getList(userID uint, offset int, limit int) ([]models.Friend, error) {
	var friends []models.Friend

	err := s.db.InnerJoins("Character").
		Where("friends.user_id = ?", userID).
		Order("friends.updated_at desc").
		Offset(offset).
		Limit(limit).
		Preload("Character.Author").
		Find(&friends).Error

	return friends, err
}