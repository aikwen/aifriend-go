package systemprompt

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
)


func (s *systemPromptStore) Update(ctx context.Context, prompt *models.SystemPrompt) error {
	return s.db.WithContext(ctx).Save(prompt).Error
}