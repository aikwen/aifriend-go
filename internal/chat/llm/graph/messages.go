package graph

import "github.com/cloudwego/eino/schema"

// FinalAssistantMessage 提取llm 回复的消息
func FinalAssistantMessage(messages []*schema.Message) *schema.Message {
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if msg == nil {
			continue
		}
		if msg.Role == schema.Assistant && len(msg.ToolCalls) == 0 {
			return msg
		}
	}
	return nil
}