package memory

import (
	"context"
	"log"
	"slices"

	"github.com/aikwen/aifriend-go/internal/chat/llm/graph"
	"github.com/aikwen/aifriend-go/internal/chat/llm/prompts"
	"github.com/aikwen/aifriend-go/internal/errs"
	"github.com/aikwen/aifriend-go/internal/store/models"
)

// Update 更新userID 名下的 friend 的记忆
func (s *memorySvc) Update(ctx context.Context, userID uint, friend *models.Friend) error{
	cnt, err := s.database.Message.CountByFriendID(ctx, friend.ID, userID)
	if err != nil {
		log.Printf("查询消息数量错误: %v", err)
		return err
	}
	// 不更新
	if cnt < 20 || cnt % 10 != 0{
		return nil
	}

	// 获取记忆提示词
	memoryMsgs, ok := s.database.Cache.SystemPrompt.Get("记忆")
	if !ok {
		memoryPromptRecords, err := s.database.SystemPrompt.GetListByTitle(ctx, "记忆")
		if err != nil {
			log.Printf("查询记忆 prompt 错误: %v", err)
			return err
		} else {
			memoryMsgs = prompts.SystemMessages(memoryPromptRecords)
			s.database.Cache.SystemPrompt.Set("记忆", memoryMsgs)
		}
	}
	memoryText := friend.Memory
	historyRecords, err := s.database.Message.GetListByFriendID(ctx, friend.ID, userID, 10, 10)
	if err != nil {
		log.Printf("查询历史最新消息 prompt 错误：%v", err)
		return err
	}
	slices.Reverse(historyRecords)
	// 历史消息
	msgs := prompts.BuildMemoryMessages(prompts.MemoryMessageInput{
		SystemMessages: memoryMsgs,
		Memory: memoryText,
		HistoryRecords: historyRecords,
	})
	// 获取总结
	outputMsgs, err := s.graph.Run(ctx, msgs, func(string){})
	if err != nil {
		log.Printf("总结memory错误：%v", err)
		return err
	}
	finalMsg := graph.FinalAssistantMessage(outputMsgs)
	if finalMsg == nil {
		log.Printf("llm 总结失败")
		return errs.ErrLLMNilFinalMessage
	}
	err = s.database.Friend.UpdateMemoryWithVersion(ctx,
		userID,
		friend.ID,
		friend.Version, finalMsg.Content)
	if err != nil {
		log.Printf("memory 更新错误：%v", err)
		return err
	}

	return nil
}