package chat

import (
	"context"
	"log"
	"slices"
	"time"

	"github.com/aikwen/aifriend-go/internal/chat/llm/graph"
	"github.com/aikwen/aifriend-go/internal/chat/llm/prompts"
	"github.com/aikwen/aifriend-go/internal/store/models"
	"github.com/aikwen/aifriend-go/pkg/task"
)

func (c *chatService) Chat(
	ctx context.Context,
	userID uint,
	friendID uint,
	message string,
) (<-chan StreamEvent, error) {
	// 查询当前用户可访问的好友及其角色信息
	friend, err := c.database.Friend.GetByIDAndUserID(ctx, userID, friendID)
	if err != nil {
		return nil, err
	}


	ch := make(chan StreamEvent, 16)

	task.Go(func() {
		defer close(ch)
		// 查询system prompt
		systemMsgs, ok := c.database.Cache.SystemPrompt.Get("回复")
		if !ok {
			systemPromptRecords, err := c.database.SystemPrompt.GetListByTitle(ctx, "回复")
			if err != nil {
				log.Printf("查询系统 prompt 错误: %v", err)
			} else {
				systemMsgs = prompts.SystemMessages(systemPromptRecords)
				c.database.Cache.SystemPrompt.Set("回复", systemMsgs)
			}
		}

		// 查询历史最新10条消息
		historyRecords, err := c.database.Message.GetLatestList(ctx, friendID, userID, 10)
		if err != nil {
			log.Printf("查询历史最新消息 prompt 错误：%v", err)
		}
		slices.Reverse(historyRecords)
		historyMsgs := prompts.HistoryMessages(historyRecords)

		// 长期记忆:
		memoryText := friend.Memory
		// 角色profile
		characterProfileText := friend.Character.Profile

		// 组装消息
		inputMessages := prompts.BuildChatMessages(prompts.ChatMessageInput{
			SystemMessages: systemMsgs,
			HistoryMessages: historyMsgs,
			Memory: memoryText,
			CharacterProfile: characterProfileText,
			UserMessage: message,
		})

		outputMsgs, err := c.graph.Run(ctx, inputMessages, func(chunk string) {
			if chunk == "" {
				return
			}
			// 往 chan 写入数据
			select {
			case ch <- StreamEvent{
				Type: EventDelta,
				Text: chunk,
			}:
			case <-ctx.Done():
				return
			}
		})

		if err != nil {
			select {
			case ch <- StreamEvent{
				Type: EventError,
				Text: err.Error(),
			}:
			case <-ctx.Done():
			}
			return
		}

		// 消息入库
		finalMsg := graph.FinalAssistantMessage(outputMsgs)
		if finalMsg == nil {
			log.Printf("llm 未回复")
			select {
			case ch <- StreamEvent{
				Type: EventError,
				Text: "llm 未返回最终回复",
			}:
			case <-ctx.Done():
			}
			return
		}

		// 将消息记录存入数据库
		if err := c.database.Message.Create(ctx, &models.Message{
			FriendID: friendID,
			UserMessage: message,
			Output: finalMsg.Content,
		}); err != nil {
			log.Printf("保存消息错误: %v", err)
		}else{
			// 更新memory
			task.Go(func() {
				memCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				if err := c.memorySvc.Update(memCtx, userID, friend); err != nil {
					log.Printf("更新记忆错误: %v", err)
				}
			})
		}
		select {
		case ch <- StreamEvent{
			Type: EventDone,
		}:
		case <-ctx.Done():
		}
	})

	return ch, nil
}