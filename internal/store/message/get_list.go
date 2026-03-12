package message

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
)


func (m *messageStore) GetRecentList(ctx context.Context, friendID uint, lastMessageID uint, userID uint, num int) ([]models.Message, error){
	var messages []models.Message

	// 确保 Character 和 Friend 都没有被删除
	query := m.db.WithContext(ctx).Table("messages").
		Select("messages.*").
		Joins("JOIN friends ON friends.id = messages.friend_id AND friends.deleted_at IS NULL").
		Joins("JOIN characters ON characters.id = friends.character_id AND characters.deleted_at IS NULL").
		// 该 friendID 必须属于当前 userID
		Where("messages.friend_id = ? AND friends.user_id = ?", friendID, userID)

	// 历史消息
	if lastMessageID > 0 {
		query = query.Where("messages.id < ?", lastMessageID)
	}

	// 寻找最新的
	err := query.Order("messages.id desc").Limit(num).Find(&messages).Error

	return messages, err
}