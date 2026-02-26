package service

import (
	"context"
	"errors"

	"github.com/aikwen/aifriend-go/internal/models"
	"github.com/aikwen/aifriend-go/internal/store"
	"github.com/aikwen/aifriend-go/pkg/storage"
	"gorm.io/gorm"
)

type UserService interface {
    GetUserInfo(ctx context.Context, userID uint) (*models.User, error)
    UpdateUserInfo(ctx context.Context, userID uint, newUsername, newProfile, newPhoto string) (*models.User, error)
}

type userService struct {
    userStore store.UserStore
    storage storage.Storage
}

func NewUserService(us store.UserStore, st storage.Storage) UserService {
    return &userService{
        userStore: us,
        storage: st,
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

func (s *userService) UpdateUserInfo(ctx context.Context, userID uint, newUsername, newProfile, newPhoto string) (*models.User, error){
    currentUser, err := s.userStore.GetByID(ctx, userID)

    if err != nil {
	    if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("用户不存在")
        }
        return nil, err
	}

    // 检查用户名是否被占用，没有占用就更新用户名
    if newUsername != currentUser.Username {
        existingUser, err := s.userStore.GetByUsername(ctx, newUsername)
        if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("检查用户名时发生数据库错误")
		}

        if existingUser != nil {
			return nil, errors.New("用户名已存在")
		}

        currentUser.Username = newUsername
    }

    // 更新用户简介
    runes := []rune(newProfile)
	if len(runes) > 500 {
		newProfile = string(runes[:500])
	}
    currentUser.Profile = newProfile

    // 更新用户头像
    if newPhoto != "" {
        if currentUser.Photo != "" && currentUser.Photo != "user/photos/default.png" {
            _ = s.storage.Delete(currentUser.Photo)
        }
        currentUser.Photo = newPhoto
	}

    if err := s.userStore.Update(ctx, currentUser); err != nil {
		return nil, errors.New("保存用户信息失败")
	}

    return currentUser, nil
}