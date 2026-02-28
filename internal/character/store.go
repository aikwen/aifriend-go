package character

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
	"gorm.io/gorm"
)

type store interface {
	create(ctx context.Context, character *models.Character) error
	getByIDAndAuthor(ctx context.Context, id uint, authorID uint) (*models.Character, error)
	update(ctx context.Context, character *models.Character) error
	delete(ctx context.Context, id uint, authorID uint) error
	getListByAuthorID(ctx context.Context, authorID uint, offset int, limit int) ([]*models.Character, error)
}

type characterStore struct {
	db *gorm.DB
}

func newCharacterStore(db *gorm.DB) store {
	return &characterStore{db: db}
}