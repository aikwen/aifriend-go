package prompts

import (
	"strings"

	"github.com/cloudwego/eino/schema"
)

type ChatMessageInput struct {
	SystemMessages   []*schema.Message
	CharacterProfile string
	Memory           string
	HistoryMessages  []*schema.Message
	UserMessage      string
}

func BuildChatMessages(in ChatMessageInput) []*schema.Message {
	messages := make([]*schema.Message, 0, len(in.SystemMessages)+len(in.HistoryMessages)+3)

	// 系统prompt
	messages = append(messages, in.SystemMessages...)

	// 角色设定
	profile := strings.TrimSpace(in.CharacterProfile)
	if profile != "" {
		messages = append(messages, schema.SystemMessage("【角色性格】\n"+ profile))
	}

	// 长期记忆
	memory := strings.TrimSpace(in.Memory)
	if memory != "" {
		messages = append(messages, schema.SystemMessage("【长期记忆】\n"+ memory))
	}

	// 历史消息
	messages = append(messages, in.HistoryMessages...)

	// 用户提问
	userMessage := strings.TrimSpace(in.UserMessage)
	if userMessage != "" {
		messages = append(messages, schema.UserMessage(userMessage))
	}

	return messages
}