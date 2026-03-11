package user

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store"
	"github.com/aikwen/aifriend-go/internal/store/models"
	"github.com/aikwen/aifriend-go/pkg/storage"
)


type Service interface {
    UpdateUserInfo(ctx context.Context, userID uint, newUsername, newProfile, newPhoto string) (*models.User, error)
    Create(ctx context.Context, user *models.User) error
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
    GetUserInfo(ctx context.Context, userID uint) (*models.User, error)
}

type userService struct {
    database *store.Database
    storage storage.Storage
}

func NewUserService(database *store.Database, st storage.Storage) Service {
    return &userService{
        database: database,
        storage: st,
    }
}