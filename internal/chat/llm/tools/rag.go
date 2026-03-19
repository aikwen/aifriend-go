package tools

import (
	"context"
	"strings"

	"github.com/aikwen/aifriend-go/internal/chat/llm/rag"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
)

const (
	SearchKnowledgeBaseToolName = "search_knowledge_base"
)

type SearchKnowledgeBaseInput struct {
	Query string `json:"query" jsonschema:"required" jsonschema_description:"要在知识库中检索的问题或关键词"`
}

type searchKnowledgeBaseTool struct {
	ragSvc rag.Service
}

func (t *searchKnowledgeBaseTool) search(ctx context.Context, input *SearchKnowledgeBaseInput) (string, error) {
	if input == nil {
		return "未提供知识库检索参数。", nil
	}

	input.Query = strings.TrimSpace(input.Query)
	if input.Query == "" {
		return "未提供有效的知识库检索问题。", nil
	}

	if t.ragSvc == nil {
		return "知识库检索暂时不可用，请不要重复调用该工具，直接基于已有上下文回答；如果缺少依据，请明确说明不知道。", nil
	}

	result, err := t.ragSvc.Search(ctx, input.Query)
	if err != nil {
		return "知识库检索暂时不可用，请不要重复调用该工具，直接基于已有上下文回答；如果缺少依据，请明确说明不知道。", nil
	}

	if strings.TrimSpace(result) == "" {
		return "未从知识库检索到相关信息。", nil
	}

	return result, nil
}

func NewSearchKnowledgeBaseTool(ragSvc rag.Service) (tool.InvokableTool, error) {
	t := &searchKnowledgeBaseTool{
		ragSvc: ragSvc,
	}

	return utils.InferTool(
		SearchKnowledgeBaseToolName,
		"当用户询问阿里云百炼平台中的事实、概念、说明文档、项目资料时，调用此工具。输入 query 为要检索的问题，输出为知识库中相关的文本片段。",
		t.search,
	)
}
