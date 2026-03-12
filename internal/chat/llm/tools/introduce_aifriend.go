package tools

import (
	"context"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

type introduceAIFriendParams struct{}


func NewIntroduceAIFriendTool() tool.InvokableTool {
	info := &schema.ToolInfo{
		Name: "introduce_aifriend",
		Desc: "返回 aifriend 网站的官方介绍。用户询问 aifriend 是什么、网站能做什么、有哪些功能、这个项目是干嘛的时，应调用此工具。",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
	}

	return utils.NewTool(info, func(ctx context.Context, input *introduceAIFriendParams) (string, error) {
		return `aifriend 是一个基于大语言模型构建的 AI 角色陪伴与互动网站。用户可以创建或选择不同的 AI 好友，与其进行流式聊天互动。平台支持角色设定、性格配置、历史消息记忆和长期记忆能力，让角色在多轮对话中表现得更真实、更连贯。系统后续也可以结合知识库检索等能力，让 AI 在聊天中提供更丰富、更准确的内容支持。`, nil
	})
}