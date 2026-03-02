package friend

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
	"gorm.io/gorm"
)


type store interface {
	getOrCreate(ctx context.Context, userID uint, characterID uint) (*models.Friend, error)
	getList(ctx context.Context, userID uint, offset int, limit int) ([]models.Friend, error)
	remove(ctx context.Context, userID uint, friendID uint) error
}

type friendStore struct{
	db *gorm.DB
}

func newFriendStore(db *gorm.DB) store{
	return &friendStore{db: db}
}

