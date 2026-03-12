package graph

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	chatmodel "github.com/aikwen/aifriend-go/internal/chat/llm/model"
	chattools "github.com/aikwen/aifriend-go/internal/chat/llm/tools"

	einomodel "github.com/cloudwego/eino/components/model"
	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

func streamAndCollect(
	ctx context.Context,
	cm einomodel.ToolCallingChatModel,
	messages []*schema.Message,
) (*schema.Message, error) {
	sr, err := cm.Stream(ctx, messages)
	if err != nil {
		return nil, err
	}
	defer sr.Close()

	var chunks []*schema.Message

	fmt.Print("assistant> ")
	for {
		chunk, err := sr.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if chunk == nil {
			continue
		}

		chunks = append(chunks, chunk)

		if chunk.Content != "" {
			fmt.Print(chunk.Content)
		}
	}
	fmt.Println()

	return schema.ConcatMessages(chunks)
}

func TestGraphLikeLoopDemo(t *testing.T) {
	ctx := context.Background()

	apiKey := os.Getenv("API_KEY")
	apiBase := os.Getenv("API_BASE")
	modelName := os.Getenv("MODEL_NAME")

	if apiKey == "" || apiBase == "" || modelName == "" {
		t.Fatal("API_KEY / API_BASE / MODEL_NAME 不能为空")
	}

	// 1. 创建模型
	cm, err := chatmodel.NewDeepseekChatModel(ctx, chatmodel.DeepseekConfig{
		APIKey:  apiKey,
		BaseURL: apiBase,
		Model:   modelName,
	})
	if err != nil {
		t.Fatalf("new deepseek model failed: %v", err)
	}

	// 2. 创建工具
	getTimeTool := chattools.NewGetTimeTool()
	introduceTool := chattools.NewIntroduceAIFriendTool()

	tools := map[string]einotool.InvokableTool{
		"get_time":            getTimeTool,
		"introduce_aifriend":  introduceTool,
	}

	// 3. 提取 ToolInfo，绑定给模型
	getTimeInfo, err := getTimeTool.Info(ctx)
	if err != nil {
		t.Fatalf("get get_time info failed: %v", err)
	}
	introduceInfo, err := introduceTool.Info(ctx)
	if err != nil {
		t.Fatalf("get introduce_aifriend info failed: %v", err)
	}

	toolModel, err := cm.WithTools([]*schema.ToolInfo{
		getTimeInfo,
		introduceInfo,
	})
	if err != nil {
		t.Fatalf("bind tools failed: %v", err)
	}

	// 4. 初始消息
	messages := []*schema.Message{
		schema.SystemMessage(
			"你是 aifriend 网站的 AI 朋友。" +
				"请根据工具的名称、描述和参数说明，自主选择合适的工具来回答问题。" +
				"当工具能够提供更准确的信息时，优先调用工具，不要直接猜测答案。",
		),
		schema.UserMessage("请先告诉我当前精确时间，然后再介绍一下 aifriend 网站可以做什么。"),
	}

	// 5. 手动实现：START -> model -> tool? -> model -> ... -> END
	const maxSteps = 10

	for step := 0; step < maxSteps; step++ {
		t.Logf("========== step %d : model ==========", step+1)

		assistantMsg, err := streamAndCollect(ctx, toolModel, messages)
		if err != nil {
			t.Fatalf("model generate failed at step %d: %v", step+1, err)
		}
		if assistantMsg == nil {
			t.Fatalf("assistant message is nil at step %d", step+1)
		}

		messages = append(messages, assistantMsg)

		t.Logf("assistant role=%q content=%q tool_calls=%+v",
			assistantMsg.Role,
			assistantMsg.Content,
			assistantMsg.ToolCalls,
		)

		// END: 没有 tool call，结束
		if len(assistantMsg.ToolCalls) == 0 {
			t.Logf("========== END at step %d ==========", step+1)
			break
		}

		// 有 tool call：执行工具
		t.Logf("========== step %d : tools ==========", step+1)

		for _, tc := range assistantMsg.ToolCalls {
			toolName := tc.Function.Name
			toolArgs := tc.Function.Arguments

			toolImpl, ok := tools[toolName]
			if !ok {
				t.Fatalf("tool %q not found", toolName)
			}

			result, err := toolImpl.InvokableRun(ctx, toolArgs)
			if err != nil {
				t.Fatalf("tool %q run failed: %v", toolName, err)
			}

			toolMsg := &schema.Message{
				Role:       schema.Tool,
				Content:    result,
				ToolCallID: tc.ID,
				Name:       toolName,
			}

			messages = append(messages, toolMsg)

			t.Logf("tool name=%q args=%q result=%q", toolName, toolArgs, result)
		}
	}

	// 6. 打印最终消息链
	t.Log("========== final messages ==========")
	for i, msg := range messages {
		t.Logf("[%d] role=%q name=%q tool_call_id=%q content=%q tool_calls=%+v",
			i, msg.Role, msg.Name, msg.ToolCallID, msg.Content, msg.ToolCalls)
	}

	// 7. 最后断言一下，确保最终一条是 assistant 正常回答
	last := messages[len(messages)-1]
	if last.Role != schema.Assistant {
		t.Fatalf("last message role = %q, want assistant", last.Role)
	}
	if last.Content == "" {
		t.Fatal("last assistant content is empty")
	}
}