package message

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
)

// GetRecentList 查询当前 userID 名下 friendID 的最近 num 条历史消息。
// 当 lastMessageID > 0 时，只返回 id 小于 lastMessageID 的更早消息。
// 消息顺序是逆序的，第一个为最新消息
func (m *messageStore) GetRecentList(ctx context.Context, friendID uint, lastMessageID uint, userID uint, num int) ([]models.Message, error){
	if num <= 0 {
		return []models.Message{}, nil
	}
	var messages []models.Message

	// 确保 Character 和 Friend 都没有被删除
	query := m.db.WithContext(ctx).Table("messages").
		Select("messages.*").
		Joins("JOIN friends ON friends.id = messages.friend_id AND friends.deleted_at IS NULL").
		Joins("JOIN characters ON characters.id = friends.character_id AND characters.deleted_at IS NULL").
		// friendID 属于当前 userID
		Where("messages.friend_id = ? AND friends.user_id = ?", friendID, userID)

	// 历史消息
	if lastMessageID > 0 {
		query = query.Where("messages.id < ?", lastMessageID)
	}

	// 寻找最新的
	err := query.Order("messages.id desc").Limit(num).Find(&messages).Error

	return messages, err
}

// GetLatestList 查询当前 userID 名下 friendID 的最近 num 条历史消息。
// 返回结果逆序，最新消息在前。
func (m *messageStore) GetLatestList(ctx context.Context, friendID uint, userID uint, num int) ([]models.Message, error) {
    return m.GetRecentList(ctx, friendID, 0, userID, num)
}


// GetListByFriendID 查询当前 userID 名下 friendID 的消息列表。
// 返回结果按消息 id 倒序排列，最新消息在前。
// offset 用于跳过最近若干条消息后，再获取指定limit数量的历史消息。
func (m *messageStore) GetListByFriendID(ctx context.Context, friendID uint, userID uint, offset int, limit int) ([]models.Message, error) {
	if limit <= 0 {
		return []models.Message{}, nil
	}

	var messages []models.Message

	err := m.db.WithContext(ctx).
		Table("messages").
		Select("messages.*").
		Joins("JOIN friends ON friends.id = messages.friend_id AND friends.deleted_at IS NULL").
		Joins("JOIN characters ON characters.id = friends.character_id AND characters.deleted_at IS NULL").
		Where("messages.friend_id = ? AND friends.user_id = ?", friendID, userID).
		Order("messages.id DESC").
		Offset(offset).
		Limit(limit).
		Find(&messages).Error
	if err != nil {
		return nil, err
	}

	return messages, nil
}