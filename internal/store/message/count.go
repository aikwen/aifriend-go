package message

import "context"

func (m *messageStore) CountByFriendID(ctx context.Context, friendID uint, userID uint) (int64, error) {
	var count int64

	err := m.db.WithContext(ctx).
		Table("messages").
		Joins("JOIN friends ON friends.id = messages.friend_id AND friends.deleted_at IS NULL").
		Joins("JOIN characters ON characters.id = friends.character_id AND characters.deleted_at IS NULL").
		Where("messages.friend_id = ? AND friends.user_id = ?", friendID, userID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}