package user

import (
	"context"
	"github.com/aikwen/aifriend-go/internal/models"
)

func (us userService) Create(ctx context.Context, user *models.User) error{
	return us.userStore.create(ctx, user)
}
