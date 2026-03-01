package friend

import (
	"github.com/aikwen/aifriend-go/internal/models"
	"gorm.io/gorm"
)


type store interface {
	getOrCreate(userID uint, characterID uint) (*models.Friend, error)
	getList(userID uint, offset int, limit int) ([]models.Friend, error)
	remove(userID uint, friendID uint) error
}

type friendStore struct{
	db *gorm.DB
}

func newFriendStore(db *gorm.DB) store{
	return &friendStore{db: db}
}

