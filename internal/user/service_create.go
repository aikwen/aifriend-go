package user

import (
	"context"
	"github.com/aikwen/aifriend-go/internal/store/models"
)

func (us *userService) Create(ctx context.Context, user *models.User) error{
	return us.database.User.Create(ctx, user)
}
