package rag

import (
	"context"
	"fmt"
	"time"

	"github.com/aikwen/aifriend-go/config"
	"github.com/aikwen/aifriend-go/internal/chat/llm/rag/embeddings"
	qdrant "github.com/qdrant/go-client/qdrant"
)

type Service interface {
	Search(ctx context.Context, query string) (string, error)
	IngestFile(ctx context.Context, filePath string) error
}

type RagConfig struct {
	Host       string
	GRPCPort   int
	APIKey     string
	Collection string
	TopK       int
	VectorSize int
}

type ragSvc struct {
	embedder   embeddings.Service
	client     *qdrant.Client
	collection string
	topK       int
}

func NewService(embedder embeddings.Service, cfg RagConfig) (Service, error) {
	if embedder == nil {
		return nil, fmt.Errorf("rag embedder is nil")
	}
	if cfg.Host == "" {
		return nil, fmt.Errorf("qdrant host is empty")
	}
	if cfg.GRPCPort <= 0 {
		return nil, fmt.Errorf("qdrant grpc port must be greater than 0")
	}
	if cfg.Collection == "" {
		return nil, fmt.Errorf("qdrant collection is empty")
	}
	if cfg.TopK <= 0 {
		return nil, fmt.Errorf("qdrant top_k must be greater than 0")
	}
	if cfg.VectorSize <= 0 {
		return nil, fmt.Errorf("qdrant vector size must be greater than 0")
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.GRPCPort)

	client, err := qdrant.NewClient(&qdrant.Config{
		Host:      cfg.Host,
		Port:      cfg.GRPCPort,
		APIKey:    cfg.APIKey,
		TLSConfig: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("new qdrant client failed, addr=%s: %w", addr, err)
	}

	ensureCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ensureCollection(ensureCtx, client, cfg.Collection, cfg.VectorSize); err != nil {
		return nil, fmt.Errorf("ensure qdrant collection failed: %w", err)
	}

	return &ragSvc{
		embedder:   embedder,
		client:     client,
		collection: cfg.Collection,
		topK:       cfg.TopK,
	}, nil
}

func NewServiceFromConfig(embedder embeddings.Service) (Service, error) {
	return NewService(embedder, RagConfig{
		Host:       config.GlobalConfig.Qdrant.Host,
		GRPCPort:   config.GlobalConfig.Qdrant.GRPCPort,
		APIKey:     config.GlobalConfig.Qdrant.APIKey,
		Collection: config.GlobalConfig.Qdrant.Collection,
		TopK:       config.GlobalConfig.Qdrant.TopK,
		VectorSize: config.GlobalConfig.LLM.EmbeddingDimensions,
	})
}