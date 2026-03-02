package friend

import (
	"context"
	"errors"

	"github.com/aikwen/aifriend-go/internal/models"
)


func (s *friendStore) remove(ctx context.Context, userID uint, friendID uint) error {
	result := s.db.WithContext(ctx).Where("id = ? AND user_id = ?", friendID, userID).Delete(&models.Friend{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("好友记录不存在或无权删除")
	}

	return nil
}