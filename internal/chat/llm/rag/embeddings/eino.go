package embeddings

import (
	"context"
	"fmt"
	"time"

	"github.com/aikwen/aifriend-go/config"
	einoembedding "github.com/cloudwego/eino/components/embedding"
	einoopenai "github.com/cloudwego/eino-ext/components/embedding/openai"
)

type EinoConfig struct {
	APIKey     string
	BaseURL    string
	Model      string
	Dimensions int
	BatchSize  int
	Timeout    time.Duration
}

type einoSvc struct {
	embedder  einoembedding.Embedder
	batchSize int
}

func NewEinoSvc(ctx context.Context, cfg EinoConfig) (Service, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("eino embedding api key is empty")
	}
	if cfg.Model == "" {
		return nil, fmt.Errorf("eino embedding model is empty")
	}
	if cfg.Dimensions <= 0 {
		return nil, fmt.Errorf("eino embedding dimensions must be greater than 0")
	}
	if cfg.BatchSize <= 0 {
		return nil, fmt.Errorf("eino embedding batch size must be greater than 0")
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}

	dimensions := cfg.Dimensions

	embedder, err := einoopenai.NewEmbedder(ctx, &einoopenai.EmbeddingConfig{
		APIKey:     cfg.APIKey,
		BaseURL:    cfg.BaseURL,
		Model:      cfg.Model,
		Dimensions: &dimensions,
		Timeout:    cfg.Timeout,
	})
	if err != nil {
		return nil, err
	}

	return &einoSvc{
		embedder:  embedder,
		batchSize: cfg.BatchSize,
	}, nil
}

func NewEinoSvcFromConfig(ctx context.Context) (Service, error) {
	llmCfg := config.GlobalConfig.LLM

	return NewEinoSvc(ctx, EinoConfig{
		APIKey:     llmCfg.APIKey,
		BaseURL:    llmCfg.APIBase,
		Model:      llmCfg.EmbeddingModel,
		Dimensions: llmCfg.EmbeddingDimensions,
		BatchSize:  llmCfg.EmbeddingBatchSize,
		Timeout:    30 * time.Second,
	})
}