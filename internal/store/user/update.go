package user

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
)

func (us *userStore) Update(ctx context.Context, user *models.User) error {
	return us.db.WithContext(ctx).Save(user).Error
}
