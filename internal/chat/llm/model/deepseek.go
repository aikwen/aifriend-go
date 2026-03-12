package model

import (
	"context"
	"fmt"

	einomodel "github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino-ext/components/model/deepseek"

	"github.com/aikwen/aifriend-go/config"
)

type DeepseekConfig struct {
	APIKey  string
	Model   string
	BaseURL string
}

func NewDeepseekChatModel(
	ctx context.Context,
	cfg DeepseekConfig,
) (einomodel.ToolCallingChatModel, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("deepseek api key is empty")
	}
	if cfg.Model == "" {
		return nil, fmt.Errorf("deepseek model is empty")
	}
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("deepseek base url is empty")
	}

	cm, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  cfg.APIKey,
		Model:   cfg.Model,
		BaseURL: cfg.BaseURL,
	})
	if err != nil {
		return nil, err
	}

	return cm, nil
}

func NewDeepseekChatModelFromConfig(
	ctx context.Context,
) (einomodel.ToolCallingChatModel, error) {
	return NewDeepseekChatModel(ctx, DeepseekConfig{
		APIKey:  config.GlobalConfig.LLM.APIKey,
		Model:   config.GlobalConfig.LLM.ModelName,
		BaseURL: config.GlobalConfig.LLM.APIBase,
	})
}