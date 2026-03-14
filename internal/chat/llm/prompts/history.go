package prompts

import (
	"strings"

	"github.com/aikwen/aifriend-go/internal/store/models"
	"github.com/cloudwego/eino/schema"
)

func HistoryMessages(records []models.Message) []*schema.Message {
	messages := make([]*schema.Message, 0, len(records)*2)

	for _, r := range records {
		userText := strings.TrimSpace(r.UserMessage)
		if userText != "" {
			messages = append(messages, schema.UserMessage(userText))
		}

		outputText := strings.TrimSpace(r.Output)
		if outputText != "" {
			messages = append(messages, schema.AssistantMessage(outputText, nil))
		}
	}

	return messages
}