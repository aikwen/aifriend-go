package prompts

import (
	"strings"

	"github.com/aikwen/aifriend-go/internal/store/models"
	"github.com/cloudwego/eino/schema"
)


type MemoryMessageInput struct {
	SystemMessages   []*schema.Message
	Memory           string
	HistoryRecords  []models.Message
}


func BuildMemoryMessages(in MemoryMessageInput) []*schema.Message {
	messages := make([]*schema.Message, 0, len(in.SystemMessages)+1)

	// 系统消息
	messages = append(messages, in.SystemMessages...)

	var builder strings.Builder

	memory := strings.TrimSpace(in.Memory)
	builder.WriteString("【原始记忆】\n")
	if memory != "" {
		builder.WriteString(memory)
	}
	builder.WriteString("\n\n【新增对话】\n")

	for _, record := range in.HistoryRecords {
		userText := strings.TrimSpace(record.UserMessage)
		if userText != "" {
			builder.WriteString("user: ")
			builder.WriteString(userText)
			builder.WriteString("\n")
		}

		outputText := strings.TrimSpace(record.Output)
		if outputText != "" {
			builder.WriteString("ai: ")
			builder.WriteString(outputText)
			builder.WriteString("\n")
		}
	}

	messages = append(messages, schema.UserMessage(builder.String()))

	return messages
}