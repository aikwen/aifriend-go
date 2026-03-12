package message

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
)


func (m *messageStore) Create(ctx context.Context, message models.Message) error {
	return m.db.WithContext(ctx).Create(message).Error
}
