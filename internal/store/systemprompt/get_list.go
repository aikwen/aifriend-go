package systemprompt

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
)


func (s *systemPromptStore) GetListByTitle(ctx context.Context, title string) ([]models.SystemPrompt, error) {
    var prompts []models.SystemPrompt
    err := s.db.WithContext(ctx).
        Where("title = ?", title).
        Order("order_number asc").
        Find(&prompts).Error
    return prompts, err
}