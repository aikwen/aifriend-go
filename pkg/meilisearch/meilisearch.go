package meilisearch

import (
	"context"
	"fmt"

	"github.com/meilisearch/meilisearch-go"
)

type Client[T comparable] interface {
	// Upsert 插入或局部更新角色数据
	Upsert(ctx context.Context, doc map[string]any) error
	UpsertInBatches(ctx context.Context, docs []map[string]any, batchSize int) error
	// SearchCharacters 搜索角色，返回匹配的 ID 列表
	Search(ctx context.Context, query string, limit, offset int) ([]T, error)
	// SetupIndex 初始化
	SetupIndex() error
	DeleteAllDocuments() error
}

type MeilisearchConfig struct {
	Host   string
	Port   int
	APIKey string
}

type meilisearchClient[T comparable] struct {
	client meilisearch.ServiceManager
	index string
	primaryKey string
}

// NewClient 初始化并返回一个 Meilisearch 客户端
// T 表示主键的类型
func NewClient[T comparable](cfg *MeilisearchConfig, index string, primaryKey string) Client[T] {
	addr := fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port)
	mc := meilisearch.New(addr, meilisearch.WithAPIKey(cfg.APIKey))
	return &meilisearchClient[T]{
		client: mc,
		index: index,
		primaryKey: primaryKey,
	}
}