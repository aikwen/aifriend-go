package character

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
)

func (s *characterService) SearchCharacters(ctx context.Context, offset int, limit int, searchQuery string) ([]*models.Character, error) {
	return s.characterStore.getListBySearchQuery(ctx, offset, limit, searchQuery)
}


