package chat

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/store/models"
)


func (c *chatService) GetHistory(ctx context.Context, friendID uint, lastMessageID uint, userID uint) ([]models.Message, error){
	return c.database.Message.GetRecentList(ctx, friendID, lastMessageID, userID, 10)
}