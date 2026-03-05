package auth

import (
	"context"
	"errors"

	"github.com/aikwen/aifriend-go/config"
	"github.com/aikwen/aifriend-go/internal/models"
	"github.com/aikwen/aifriend-go/pkg/auth"
	"gorm.io/gorm"
)

var ErrUserAlreadyExists = errors.New("用户名已存在")

func (as *authService) Register(ctx context.Context, username, password string, jwtConf *config.JWTConfig) (*models.User, string, string, error) {
    _, err := as.userService.GetByUsername(ctx, username)
    if err == nil {
        return nil, "", "", ErrUserAlreadyExists
    }

    if !errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, "", "", err
    }

    hashedPassword, err := auth.HashPassword(password)
    if err != nil {
        return nil, "", "", errors.New("密码加密失败")
    }

    newUser := &models.User{
        Username: username,
        Password: hashedPassword,
    }
    if err := as.userService.Create(ctx, newUser); err != nil {
        return nil, "", "", err
    }

    accessToken, err := auth.GenerateAccessToken(newUser.Username, newUser.ID, []byte(jwtConf.AccessSecret))
    if err != nil {
        return nil, "", "", err
    }
    refreshToken, err := auth.GenerateRefreshToken(newUser.Username, newUser.ID, []byte(jwtConf.RefreshSecret))
    if err != nil {
        return nil, "", "", err
    }

    return newUser, accessToken, refreshToken, nil
}