package user


import (
	"context"
	"errors"

	"github.com/aikwen/aifriend-go/internal/models"
	"gorm.io/gorm"
)


func (s *userService) GetUserInfo(ctx context.Context, userID uint) (*models.User, error) {
    user, err := s.userStore.getByID(ctx, userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("用户不存在")
        }
        return nil, err
    }
    return user, nil
}


func (s *userService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
    return s.userStore.getByUsername(ctx, username)
}


func (s *userService) GetByID(ctx context.Context, id uint) (*models.User, error) {
    return s.userStore.getByID(ctx, id)
}
