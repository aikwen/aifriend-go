package embeddings

import (
	"fmt"

	"github.com/aikwen/aifriend-go/config"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type OpenAIConfig struct {
	APIKey     string
	BaseURL    string
	Model      string
	Dimensions int
	BatchSize  int
}

type openaiSvc struct {
	client     openai.Client
	model      string
	dimensions int64
	batchSize  int
}

func NewOpenAISvc(cfg OpenAIConfig) (Service, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("openai api key is empty")
	}
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("openai base url is empty")
	}
	if cfg.Model == "" {
		return nil, fmt.Errorf("embedding model is empty")
	}
	if cfg.Dimensions <= 0 {
		return nil, fmt.Errorf("embedding dimensions must be greater than 0")
	}
	if cfg.BatchSize <= 0 {
		return nil, fmt.Errorf("embedding batch size must be greater than 0")
	}

	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseURL),
	)

	return &openaiSvc{
		client:     client,
		model:      cfg.Model,
		dimensions: int64(cfg.Dimensions),
		batchSize:  cfg.BatchSize,
	}, nil
}

func NewOpenAISvcFromConfig() (Service, error) {
	llmCfg := config.GlobalConfig.LLM

	return NewOpenAISvc(OpenAIConfig{
		APIKey:     llmCfg.APIKey,
		BaseURL:    llmCfg.APIBase,
		Model:      llmCfg.EmbeddingModel,
		Dimensions: llmCfg.EmbeddingDimensions,
		BatchSize:  llmCfg.EmbeddingBatchSize,
	})
}