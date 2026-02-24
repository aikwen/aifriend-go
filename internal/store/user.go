package store

import (
	"context"
	"errors"

	"github.com/aikwen/aifriend-go/internal/models"
	"gorm.io/gorm"
)

type UserStore interface {
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
}


type userStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) UserStore {
	return &userStore{
		db: db,
	}
}

func (us *userStore) Create(ctx context.Context, user *models.User) error {
	return us.db.WithContext(ctx).Create(user).Error
}

func (us *userStore) Update(ctx context.Context, user *models.User) error {
	return us.db.WithContext(ctx).Save(user).Error
}

func (us *userStore) GetByUsername(ctx context.Context, username string) (*models.User, error){
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

func (us *userStore) GetByID(ctx context.Context, id uint) (*models.User, error){
	var user models.User
    err := us.db.WithContext(ctx).First(&user, id).Error
    return &user, err
}