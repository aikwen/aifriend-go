 package character

 import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
)


// GetCharacterList 获取虚拟角色列表
func (s *characterService) GetCharacterList(ctx context.Context, authorID uint, offset int) ([]*models.Character, error) {
	return s.characterStore.getListByAuthorID(ctx, authorID, offset, 20)
}