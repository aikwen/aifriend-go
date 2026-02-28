package auth

import (
	"context"
	"errors"

	"github.com/aikwen/aifriend-go/internal/models"
	"github.com/aikwen/aifriend-go/pkg/auth"
	"gorm.io/gorm"
)


func (as *authService) Login(ctx context.Context, username, password string) (*models.User, string, string, error){
    user, err := as.userService.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", errors.New("用户名或密码错误")
		}
		return nil, "", "", err
	}

    isMatch, err := auth.CheckPassword(password, user.Password)
    if err != nil {
        return nil, "","", errors.New("系统校验异常，请稍后再试")
    }
    if !isMatch {
		return nil, "", "", errors.New("用户名或密码错误")
	}

    accessToken, err := auth.GenerateAccessToken(user.Username, user.ID, []byte(as.jwtConf.AccessSecret))
    if err != nil {
		return nil, "", "", errors.New("颁发 AccessToken 失败")
	}

    refreshToken, err := auth.GenerateRefreshToken(user.Username, user.ID, []byte(as.jwtConf.RefreshSecret))
	if err != nil {
		return nil, "", "", errors.New("颁发 RefreshToken 失败")
	}

    return user, accessToken, refreshToken, nil
}