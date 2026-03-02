package character

import (
	"context"

	"github.com/aikwen/aifriend-go/internal/models"
)


func (cs *characterStore) exist(ctx context.Context, characterId uint) (bool, error){
	var count int64
	err := cs.db.Model(&models.Character{}).WithContext(ctx).Where("id = ?", characterId).Count(&count).Error
	if err != nil {
		return false, err // 数据库发生其他查询错误
	}

	return count > 0, nil
}