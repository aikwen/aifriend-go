package systemprompt

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
	"gorm.io/gorm"
)


type Store interface {
	GetListByTitle(ctx context.Context, title string) ([]models.SystemPrompt, error)
	Create(ctx context.Context, prompt *models.SystemPrompt) error
	Update(ctx context.Context, prompt *models.SystemPrompt) error
}


type systemPromptStore struct {
    db *gorm.DB
}

func NewSystemPromptStore(db *gorm.DB) Store {
    return &systemPromptStore{db: db}
}