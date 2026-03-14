package character

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
	"gorm.io/gorm"
)

type Store interface {
	Create(ctx context.Context, character *models.Character) error
	GetByIDAndAuthor(ctx context.Context, id uint, authorID uint) (*models.Character, error)
	Update(ctx context.Context, character *models.Character) error
	Delete(ctx context.Context, id uint, authorID uint) error
	GetListByAuthorID(ctx context.Context, authorID uint, offset int, limit int) ([]*models.Character, error)
	GetListBySearchQuery(ctx context.Context, offset int, limit int, searchQuery string) ([]*models.Character, error)
	Exist(ctx context.Context, characterID uint) (bool, error)
	GetByID(ctx context.Context, id uint) (*models.Character, error)
}

type characterStore struct {
	db *gorm.DB
}

func NewCharacterStore(db *gorm.DB) Store {
	return &characterStore{db: db}
}
