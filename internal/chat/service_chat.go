package chat

import (
	"context"

	"github.com/cloudwego/eino/schema"
)

func (c *chatService) Chat(
	ctx context.Context,
	userID uint,
	friendID uint,
	message string,
) (<-chan StreamEvent, error) {
	_ = userID
	_ = friendID

	ch := make(chan StreamEvent, 16)

	go func() {
		defer close(ch)

		messages := []*schema.Message{
			schema.SystemMessage(
				"你是 aifriend 网站中的 AI 朋友。" +
					"你的首要任务是自然、准确地回答用户当前的问题。" +
					"当工具能够提供更准确的信息时，优先调用工具，不要直接猜测答案。" +
					"当用户的问题不需要调用工具时，就直接自然地回答用户。" +
					"只有当用户明确询问 aifriend 网站、功能、项目介绍时，才介绍 aifriend。" +
					"不要在普通问题回答后主动引导到网站功能或项目介绍。",
			),
			schema.UserMessage(message),
		}

		_, err := c.graph.Run(ctx, messages, func(chunk string) {
			if chunk == "" {
				return
			}
			ch <- StreamEvent{
				Type: EventDelta,
				Text: chunk,
			}
		})
		if err != nil {
			ch <- StreamEvent{
				Type: EventError,
				Text: err.Error(),
			}
			return
		}

		ch <- StreamEvent{
			Type: EventDone,
		}
	}()

	return ch, nil
}