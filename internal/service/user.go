package service

import (
    "context"
    "errors"

    "github.com/aikwen/aifriend-go/internal/models"
    "github.com/aikwen/aifriend-go/internal/store"
    "gorm.io/gorm"
)

type UserService interface {
    GetUserInfo(ctx context.Context, userID uint) (*models.User, error)
}

type userService struct {
    userStore store.UserStore
}

func NewUserService(us store.UserStore) UserService {
    return &userService{
        userStore: us,
    }
}

func (s *userService) GetUserInfo(ctx context.Context, userID uint) (*models.User, error) {
    user, err := s.userStore.GetByID(ctx, userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("用户不存在")
        }
        return nil, err
    }
    return user, nil
}