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
	Exists(ctx context.Context, userID uint, friendID uint) (bool, error)
}

type friendStore struct {
	db *gorm.DB
}

func NewFriendStore(db *gorm.DB) Store {
	return &friendStore{db: db}
}
