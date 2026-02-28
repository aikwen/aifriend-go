package user

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
	"gorm.io/gorm"
)

type store interface {
	create(ctx context.Context, user *models.User) error
	update(ctx context.Context, user *models.User) error
	getByUsername(ctx context.Context, username string) (*models.User, error)
	getByID(ctx context.Context, id uint) (*models.User, error)
}


type userStore struct {
	db *gorm.DB
}

func newUserStore(db *gorm.DB) store {
	return &userStore{
		db: db,
	}
}