package chat

import (
	"context"
	"fmt"

	"github.com/aikwen/aifriend-go/internal/chat/llm/graph"
	"github.com/aikwen/aifriend-go/internal/chat/llm/memory"
	chatmodel "github.com/aikwen/aifriend-go/internal/chat/llm/model"
	"github.com/aikwen/aifriend-go/internal/chat/llm/rag"
	"github.com/aikwen/aifriend-go/internal/chat/llm/rag/embeddings"
	chattools "github.com/aikwen/aifriend-go/internal/chat/llm/tools"
	"github.com/aikwen/aifriend-go/internal/store"
	"github.com/aikwen/aifriend-go/internal/store/models"

	einotool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type Service interface {
	Chat(ctx context.Context, userID uint, friendID uint, message string) (<-chan StreamEvent, error)
	GetHistory(ctx context.Context, friendID uint, lastMessageID uint, userID uint) ([]models.Message, error)
}

type chatService struct {
	database  *store.Database
	graph     *graph.Graph
	memorySvc memory.Service
}

func NewChatService(database *store.Database) (Service, error) {
	ctx := context.Background()
	embeddingSvc, err := embeddings.NewEinoSvcFromConfig(ctx)
	if err != nil {
		return nil, err
	}

	ragSvc, err := rag.NewServiceFromConfig(embeddingSvc)
	if err != nil {
		return nil, err
	}

	// 创建工具
	getTimeTool := chattools.NewGetTimeTool()
	introduceTool := chattools.NewIntroduceAIFriendTool()
	ragTool, err := chattools.NewSearchKnowledgeBaseTool(ragSvc)
	if err != nil {
		return nil, err
	}

	tools := map[string]einotool.InvokableTool{
		chattools.TimeToolName:                getTimeTool,
		chattools.IntroduceAIFriendToolName:   introduceTool,
		chattools.SearchKnowledgeBaseToolName: ragTool,
	}

	// 提取 ToolInfo
	getTimeInfo, err := getTimeTool.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("get get_time tool info failed: %w", err)
	}

	introduceInfo, err := introduceTool.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("get introduce_aifriend tool info failed: %w", err)
	}

	ragInfo, err := ragTool.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("get search_knowledge_base tool info failed: %w", err)
	}

	// 创建基础模型
	cm, err := chatmodel.NewDeepseekChatModelFromConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("init deepseek chat model failed: %w", err)
	}

	// 给模型绑定 tools
	toolModel, err := cm.WithTools([]*schema.ToolInfo{
		getTimeInfo,
		introduceInfo,
		ragInfo,
	})

	if err != nil {
		return nil, fmt.Errorf("bind tools to model failed: %w", err)
	}

	// 创建 graph
	g := graph.NewGraph(toolModel, tools)

	return &chatService{
		database:  database,
		graph:     g,
		memorySvc: memory.NewService(g, database),
	}, nil
}
