package auth



import (
	"context"
	"errors"
	"strconv"

	"github.com/aikwen/aifriend-go/pkg/auth"
)


func (as *authService) RefreshToken(ctx context.Context, refreshTokenString string) (string, string, error){
    claims, err := auth.ParseRefreshToken(refreshTokenString, []byte(as.jwtConf.RefreshSecret))
    if err != nil {
        return "", "", errors.New("无效的刷新令牌或已过期，请重新登录")
    }

    userID64, err := strconv.ParseUint(claims.Subject, 10, 64)
    if err != nil {
        return "", "", errors.New("令牌数据异常")
    }
    userID := uint(userID64)

    user, err := as.userService.GetByID(ctx, userID)
	if err != nil {
		return "", "", errors.New("用户不存在或状态异常")
	}

    newAccessToken, err := auth.GenerateAccessToken(user.Username, user.ID, []byte(as.jwtConf.AccessSecret))
	if err != nil {
		return "", "", errors.New("生成新的 AccessToken 失败")
	}

    newRefreshToken := refreshTokenString
    if as.jwtConf.RotateRefreshTokens {
        newRefreshToken, err = auth.GenerateRefreshToken(user.Username, user.ID, []byte(as.jwtConf.RefreshSecret))
		if err != nil {
			return "", "", errors.New("生成新的 RefreshToken 失败")
		}
    }

    return newAccessToken, newRefreshToken, nil
}