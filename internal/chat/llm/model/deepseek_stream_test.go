package model

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/cloudwego/eino/schema"
)

func TestDeepseekStreamToolCall(t *testing.T) {
	ctx := context.Background()

	apiKey := os.Getenv("API_KEY")
	apiBase := os.Getenv("API_BASE")
	modelName := os.Getenv("MODEL_NAME")

	if apiKey == "" || apiBase == "" || modelName == "" {
		t.Fatal("API_KEY / API_BASE / MODEL_NAME 不能为空")
	}

	cm, err := NewDeepseekChatModel(ctx, DeepseekConfig{
		APIKey:  apiKey,
		BaseURL: apiBase,
		Model:   modelName,
	})
	if err != nil {
		t.Fatalf("new deepseek model failed: %v", err)
	}

	toolModel, err := cm.WithTools([]*schema.ToolInfo{
		{
			Name: "get_time",
			Desc: "返回当前精确时间。用户询问现在几点、当前时间、日期时间时，应调用此工具。",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
		},
	})
	if err != nil {
		t.Fatalf("with tools failed: %v", err)
	}

	sr, err := toolModel.Stream(ctx, []*schema.Message{
		schema.SystemMessage(
			"当工具能够提供更准确的信息时，优先调用工具，不要直接猜测答案。" +
				"如果需要调用工具，请直接进行工具调用，不要先输出解释性文本。",
		),
		schema.UserMessage("帮我查一下现在的精确时间。"),
	})
	if err != nil {
		t.Fatalf("stream failed: %v", err)
	}
	defer sr.Close()

	i := 0
	for {
		chunk, err := sr.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("recv failed: %v", err)
		}

		t.Logf(
			"chunk[%d] role=%q content=%q tool_calls=%+v response_meta=%+v",
			i,
			chunk.Role,
			chunk.Content,
			chunk.ToolCalls,
			chunk.ResponseMeta,
		)
		i++
	}

	if i == 0 {
		t.Fatal("no stream chunks received")
	}
}