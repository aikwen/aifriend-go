 package character

 import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
)

// GetCharacter 根据角色 ID 获取单个AI角色的详细信息
func (s *characterService) GetCharacter(ctx context.Context, id uint, authorID uint) (*models.Character, error) {
	return s.characterStore.getByIDAndAuthor(ctx, id, authorID)
}