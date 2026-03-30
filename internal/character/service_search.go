package character

import (
	"context"
	"log"
	"strings"

	"github.com/aikwen/aifriend-go/internal/store/models"
)

func (s *characterService) SearchCharacters(ctx context.Context, offset int, limit int, searchQuery string) ([]*models.Character, error) {
	searchQuery = strings.TrimSpace(searchQuery)

	if searchQuery == "" {
		return s.database.Character.GetListBySearchQuery(ctx, offset, limit, "")
	}

	ids, err := s.msClient.Search(ctx, searchQuery, limit, offset)

	if err != nil {
		log.Printf("[搜索降级] Meilisearch 搜索异常 (%v)，正在切回 MySQL 原生检索...", err)
		// 自动降级
		return s.database.Character.GetListBySearchQuery(ctx, offset, limit, searchQuery)
	}

	log.Printf("关键字: '%s', 命中数量: %d, ID列表: %v", searchQuery, len(ids), ids)


	if len(ids) == 0 {
		return []*models.Character{}, nil
	}

	return s.database.Character.GetListByIDsWithOrder(ctx, ids)
}


