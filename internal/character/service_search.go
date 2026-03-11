package character

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
)

func (s *characterService) SearchCharacters(ctx context.Context, offset int, limit int, searchQuery string) ([]*models.Character, error) {
	return s.database.Character.GetListBySearchQuery(ctx, offset, limit, searchQuery)
}


