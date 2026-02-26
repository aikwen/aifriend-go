package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/aikwen/aifriend-go/internal/models"
	"github.com/aikwen/aifriend-go/internal/store"
	"github.com/aikwen/aifriend-go/pkg/auth"
	"gorm.io/gorm"
)
type AuthService interface {
    // accessToken, refreshToken
    Login(ctx context.Context, username, password string) (*models.User, string, string, error)
    Register(ctx context.Context, username, password string) (*models.User, string, string, error)

    RefreshToken(ctx context.Context, refreshTokenString string) (string, string, error)
}

type authService struct {
    userStore     store.UserStore
    accessSecret  []byte
    refreshSecret []byte
    rotateRefreshTokens bool
}


func NewAuthService(us store.UserStore, accessSecret string, refreshSecret string, rotate bool) AuthService {
    return &authService{
        userStore:     us,
        accessSecret:  []byte(accessSecret),
        refreshSecret: []byte(refreshSecret),
        rotateRefreshTokens: rotate,
    }
}

func (as *authService) Login(ctx context.Context, username, password string) (*models.User, string, string, error){
    user, err := as.userStore.GetByUsername(ctx, username)
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

    accessToken, err := auth.GenerateAccessToken(user.Username, user.ID, as.accessSecret)
    if err != nil {
		return nil, "", "", errors.New("颁发 AccessToken 失败")
	}

    refreshToken, err := auth.GenerateRefreshToken(user.Username, user.ID, as.refreshSecret)
	if err != nil {
		return nil, "", "", errors.New("颁发 RefreshToken 失败")
	}

    return user, accessToken, refreshToken, nil
}

func (as *authService) Register(ctx context.Context, username, password string) (*models.User, string, string, error) {
    _, err := as.userStore.GetByUsername(ctx, username)
    if err == nil {
        return nil, "", "", errors.New("用户名已存在")
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
    if err := as.userStore.Create(ctx, newUser); err != nil {
        return nil, "", "", err
    }

    accessToken, err := auth.GenerateAccessToken(newUser.Username, newUser.ID, as.accessSecret)
    if err != nil {
        return nil, "", "", err
    }
    refreshToken, err := auth.GenerateRefreshToken(newUser.Username, newUser.ID, as.refreshSecret)
    if err != nil {
        return nil, "", "", err
    }

    return newUser, accessToken, refreshToken, nil
}

func (as *authService) RefreshToken(ctx context.Context, refreshTokenString string) (string, string, error){
    claims, err := auth.ParseRefreshToken(refreshTokenString, as.refreshSecret)
    if err != nil {
        return "", "", errors.New("无效的刷新令牌或已过期，请重新登录")
    }

    userID64, err := strconv.ParseUint(claims.Subject, 10, 64)
    if err != nil {
        return "", "", errors.New("令牌数据异常")
    }
    userID := uint(userID64)

    user, err := as.userStore.GetByID(ctx, userID)
	if err != nil {
		return "", "", errors.New("用户不存在或状态异常")
	}

    newAccessToken, err := auth.GenerateAccessToken(user.Username, user.ID, as.accessSecret)
	if err != nil {
		return "", "", errors.New("生成新的 AccessToken 失败")
	}

    newRefreshToken := refreshTokenString
    if as.rotateRefreshTokens {
        newRefreshToken, err = auth.GenerateRefreshToken(user.Username, user.ID, as.refreshSecret)
		if err != nil {
			return "", "", errors.New("生成新的 RefreshToken 失败")
		}
    }

    return newAccessToken, newRefreshToken, nil
}