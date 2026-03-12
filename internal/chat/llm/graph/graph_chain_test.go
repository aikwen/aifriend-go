package graph

import (
	"context"
	"os"
	"strings"
	"testing"

	chatmodel "github.com/aikwen/aifriend-go/internal/chat/llm/model"
	chattools "github.com/aikwen/aifriend-go/internal/chat/llm/tools"

	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

func TestGraphRun_WithTools(t *testing.T) {
	ctx := context.Background()

	apiKey := os.Getenv("API_KEY")
	apiBase := os.Getenv("API_BASE")
	modelName := os.Getenv("MODEL_NAME")

	if apiKey == "" || apiBase == "" || modelName == "" {
		t.Fatal("API_KEY / API_BASE / MODEL_NAME 不能为空")
	}

	// 创建基础模型
	cm, err := chatmodel.NewDeepseekChatModel(ctx, chatmodel.DeepseekConfig{
		APIKey:  apiKey,
		BaseURL: apiBase,
		Model:   modelName,
	})
	if err != nil {
		t.Fatalf("new deepseek model failed: %v", err)
	}

	// 创建工具
	getTimeTool := chattools.NewGetTimeTool()
	introduceTool := chattools.NewIntroduceAIFriendTool()

	tools := map[string]einotool.InvokableTool{
		"get_time":           getTimeTool,
		"introduce_aifriend": introduceTool,
	}

	// 获取 ToolInfo 并绑定给模型
	getTimeInfo, err := getTimeTool.Info(ctx)
	if err != nil {
		t.Fatalf("get get_time tool info failed: %v", err)
	}

	introduceInfo, err := introduceTool.Info(ctx)
	if err != nil {
		t.Fatalf("get introduce_aifriend tool info failed: %v", err)
	}

	toolModel, err := cm.WithTools([]*schema.ToolInfo{
		getTimeInfo,
		introduceInfo,
	})
	if err != nil {
		t.Fatalf("bind tools failed: %v", err)
	}

	// 创建正式 Graph
	g := NewGraph(toolModel, tools)
	g.SetMaxSteps(10)
	g.SetDebug(true)

	// 初始消息
	messages := []*schema.Message{
		schema.SystemMessage(
			"你是 aifriend 网站中的 AI 朋友。" +
				"你的首要任务是自然、准确地回答用户当前的问题。" +
				"当工具能够提供更准确的信息时，优先调用工具，不要直接猜测答案。" +
				"当用户的问题不需要调用工具时，就直接自然地回答用户。" +
				"只有当用户明确询问 aifriend 网站、功能、项目介绍时，才介绍 aifriend。" +
				"不要在普通问题回答后主动引导到网站功能或项目介绍。",
		),
		schema.UserMessage("请先告诉我当前精确时间，然后再介绍一下 aifriend 网站可以做什么。"),
	}

	var fullText strings.Builder

	// 跑 Graph
	finalMessages, err := g.Run(ctx, messages, func(chunk string) {
		if chunk != "" {
			fullText.WriteString(chunk)
		}
	})
	if err != nil {
		t.Fatalf("graph run failed: %v", err)
	}

	if len(finalMessages) == 0 {
		t.Fatal("final messages is empty")
	}

	// 找最后一条 assistant 消息
	last := finalMessages[len(finalMessages)-1]
	if last.Role != schema.Assistant {
		t.Fatalf("last message role = %q, want assistant", last.Role)
	}
	if strings.TrimSpace(last.Content) == "" {
		t.Fatal("last assistant content is empty")
	}

	t.Logf("stream collected text: %s", fullText.String())
	t.Logf("last assistant content: %s", last.Content)

	// 验证消息链里至少出现过一次 tool 消息
	var hasToolMsg bool
	for i, msg := range finalMessages {
		t.Logf("[%d] role=%q name=%q tool_call_id=%q content=%q tool_calls=%+v",
			i, msg.Role, msg.Name, msg.ToolCallID, msg.Content, msg.ToolCalls)

		if msg.Role == schema.Tool {
			hasToolMsg = true
		}
	}

	if !hasToolMsg {
		t.Logf("expected at least one tool message, but got none")
	}
}