package message

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
	"gorm.io/gorm"
)


type Store interface {
	Create(ctx context.Context, message *models.Message) error
	GetRecentList(ctx context.Context, friendID uint, lastMessageID uint, userID uint, num int) ([]models.Message, error)
	GetLatestList(ctx context.Context, friendID uint, userID uint, num int) ([]models.Message, error)
}


type messageStore struct{
	db *gorm.DB
}

func NewMessageStore(db *gorm.DB) Store {
	return &messageStore{
		db:db,
	}
}