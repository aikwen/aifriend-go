package user

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
)

// Create 创建新用户
func (us *userStore) Create(ctx context.Context, user *models.User) error {
	return us.db.WithContext(ctx).Create(user).Error
}
