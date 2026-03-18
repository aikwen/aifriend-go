package friend

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
	"gorm.io/gorm"
)

type Store interface {
	GetOrCreate(ctx context.Context, userID uint, characterID uint) (*models.Friend, error)
	GetList(ctx context.Context, userID uint, offset int, limit int) ([]models.Friend, error)
	Remove(ctx context.Context, userID uint, friendID uint) error
	GetByIDAndUserID(ctx context.Context, userID uint, friendID uint) (*models.Friend, error)
	UpdateMemoryWithVersion(ctx context.Context, userID uint, friendID uint, oldVersion uint, memory string) error
}

type friendStore struct {
	db *gorm.DB
}

func NewFriendStore(db *gorm.DB) Store {
	return &friendStore{db: db}
}
