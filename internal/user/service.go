package user

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
	"github.com/aikwen/aifriend-go/pkg/storage"
	"gorm.io/gorm"
)


type Service interface {
    GetUserInfo(ctx context.Context, userID uint) (*models.User, error)
    UpdateUserInfo(ctx context.Context, userID uint, newUsername, newProfile, newPhoto string) (*models.User, error)
    Create(ctx context.Context, user *models.User) error
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
}

type userService struct {
    userStore store
    storage storage.Storage
}

func NewUserService(db *gorm.DB, st storage.Storage) Service {
    return &userService{
        userStore: newUserStore(db),
        storage: st,
    }
}