package meilisearch

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/meilisearch/meilisearch-go"
)

// Search 搜索角色，返回匹配的 ID 列表
func (c *meilisearchClient[T]) Search(ctx context.Context, query string, limit, offset int) ([]T, error) {
	req := &meilisearch.SearchRequest{
		Limit:  int64(limit),
		Offset: int64(offset),
		Filter: []string{"deleted_at IS NULL OR deleted_at NOT EXISTS"},
		Sort:   []string{fmt.Sprintf("%s:desc", c.primaryKey)},
		AttributesToRetrieve: []string{c.primaryKey},
	}

	res, err := c.client.Index(c.index).SearchWithContext(ctx, query, req)
	if err != nil {
		log.Printf("[Meilisearch] 索引 %s 搜索失败: %v\n", c.index, err)
		return nil, err
	}

	if len(res.Hits) == 0 {
		return []T{}, nil
	}


	var extractors []map[string]T

	hitsBytes, err := json.Marshal(res.Hits)
	if err != nil {
		log.Printf("[Meilisearch] 解析 Hits 失败: %v\n", err)
		return nil, err
	}

	if err := json.Unmarshal(hitsBytes, &extractors); err != nil {
		log.Printf("[Meilisearch] 反序列化 ID 字典失败: %v\n", err)
		return nil, err
	}

	var ids []T
	for _, ext := range extractors {
		if val, ok := ext[c.primaryKey]; ok {
			ids = append(ids, val)
		}
	}

	return ids, nil
}