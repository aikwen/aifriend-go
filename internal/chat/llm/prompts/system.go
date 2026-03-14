package prompts

import (
	"strings"

	"github.com/aikwen/aifriend-go/internal/store/models"
	"github.com/cloudwego/eino/schema"
)

func SystemMessages(prompts []models.SystemPrompt) []*schema.Message {
	messages := make([]*schema.Message, 0, len(prompts))

	for _, p := range prompts {
		content := strings.TrimSpace(p.Prompt)
		if content == "" {
			continue
		}
		messages = append(messages, schema.SystemMessage(content))
	}

	return messages
}