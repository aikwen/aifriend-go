package user

import (
	"context"
	"errors"

	"github.com/aikwen/aifriend-go/internal/models"
	"gorm.io/gorm"
)

// getByUsername 根据用户名查询
func (us *userStore) getByUsername(ctx context.Context, username string) (*models.User, error){
	var user models.User
	err := us.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		// 没有找到记录
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err // 明确返回找不到
		}
		// 数据库连接断了等严重错误
		return nil, err
	}

	return &user, nil
}

// getByID 根据用户id查询
func (us *userStore) getByID(ctx context.Context, id uint) (*models.User, error){
	var user models.User
    err := us.db.WithContext(ctx).First(&user, id).Error
    return &user, err
}