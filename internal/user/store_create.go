package user

import (
	"context"
	"github.com/aikwen/aifriend-go/internal/models"
)

// create 创建新用户
func (us *userStore) create(ctx context.Context, user *models.User) error {
	return us.db.WithContext(ctx).Create(user).Error
}