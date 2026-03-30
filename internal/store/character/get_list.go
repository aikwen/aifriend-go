package character

import (
	"context"
	"fmt"
	"strings"

	"github.com/aikwen/aifriend-go/internal/store/models"
	"gorm.io/gorm"
)

// GetListByAuthorID 根据用户id和 offset 获取 limit 条数据
func (s *characterStore) GetListByAuthorID(ctx context.Context, authorID uint, offset int, limit int) ([]*models.Character, error) {
	var characters []*models.Character
	err := s.db.WithContext(ctx).
		Where("author_id = ?", authorID).
		Preload("Author").
		Order("id desc").
		Offset(offset).
		Limit(limit).
		Find(&characters).Error
	return characters, err
}

// getListbySearchQuery 根据关键字搜索虚拟角色
func (s *characterStore) GetListBySearchQuery(ctx context.Context, offset int, limit int, searchQuery string) ([]*models.Character, error) {
	var characters []*models.Character

	searchQuery = strings.TrimSpace(searchQuery)

	query := s.db.WithContext(ctx).Preload("Author")

	if searchQuery != "" {
		query = query.Where("MATCH (name, profile) AGAINST (? IN NATURAL LANGUAGE MODE)", searchQuery).
                  Order(gorm.Expr("MATCH (name, profile) AGAINST (? IN NATURAL LANGUAGE MODE) DESC", searchQuery)).
				  Order("id DESC")
	} else {
		query = query.Order("id DESC")
	}

	err := query.
		Offset(offset).
		Limit(limit).
		Find(&characters).
		Error

	return characters, err
}


// GetListByIDsWithOrder 根据指定的 ID 列表获取角色，并严格按照 ID 列表的顺序返回
func (s *characterStore) GetListByIDsWithOrder(ctx context.Context, ids []uint) ([]*models.Character, error) {
	if len(ids) == 0 {
		return []*models.Character{}, nil
	}

	var characters []*models.Character

	// 构建 MySQL 的 FIELD 排序语句
	// 传进来是 [5, 2, 8]，那么生成的 SQL 就是: ORDER BY FIELD(id, 5, 2, 8)
	args := make([]any, len(ids))
	placeholders := make([]string, len(ids))
	for i, id := range ids {
		args[i] = id
		placeholders[i] = "?"
	}
	orderSQL := fmt.Sprintf("FIELD(id, %s)", strings.Join(placeholders, ","))

	err := s.db.WithContext(ctx).
		Preload("Author").
		Where("id IN ?", ids).
		Order(gorm.Expr(orderSQL, args...)).
		Find(&characters).Error

	return characters, err
}