package user

import (
	"context"
	"github.com/aikwen/aifriend-go/internal/models"
)

func (us *userStore) update(ctx context.Context, user *models.User) error {
	return us.db.WithContext(ctx).Save(user).Error
}