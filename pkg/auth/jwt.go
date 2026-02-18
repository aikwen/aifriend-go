package auth

import (
	"time"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

const (
	Issuer = "AiFriend"

	KeyID = "v1"

	AccessTokenAudienceName = "user.access-token"
	RefreshTokenAudienceName = "user.refresh-token"

	AccessTokenLifetime = 2 * time.Hour
	RefreshTokenLifetime = 7 * 24 * time.Hour

	RotateRefreshTokens = true
)

type ClaimsMessage struct {
	Name string `json:"name"` // Username
	jwt.RegisteredClaims
}

func GenerateAccessToken(username string, userID uint, secret []byte) (string, error) {
	return generateToken(username, userID, AccessTokenAudienceName, AccessTokenLifetime, secret)
}

func GenerateRefreshToken(username string, userID uint, secret []byte) (string, error) {
	return generateToken(username, userID, RefreshTokenAudienceName, RefreshTokenLifetime, secret)
}

func generateToken(username string, userID uint, audience string, duration time.Duration, secret []byte) (string, error) {
	registeredClaims := jwt.RegisteredClaims{
		Issuer:   Issuer,
		Audience: jwt.ClaimStrings{audience},
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Subject:  fmt.Sprint(userID),
	}

	registeredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(duration))


	// Declare the token with the HS256 algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &ClaimsMessage{
		Name:             username,
		RegisteredClaims: registeredClaims,
	})
	token.Header["kid"] = KeyID

	// Create the JWT string.
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseAccessToken(tokenString string, secret []byte) (*ClaimsMessage, error) {
	return parseToken(tokenString, secret, AccessTokenAudienceName)
}

func ParseRefreshToken(tokenString string, secret []byte) (*ClaimsMessage, error) {
	return parseToken(tokenString, secret, RefreshTokenAudienceName)
}

func parseToken(tokenString string, secret []byte, expectedAudience string) (*ClaimsMessage, error) {
	if tokenString == "" {
		return nil, fmt.Errorf("token 校验失败")
	}

	claims := &ClaimsMessage{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, fmt.Errorf("非法的签名算法: %v", t.Header["alg"])
		}

		if kid, ok := t.Header["kid"].(string); ok {
			if kid == "v1" {
				return secret, nil
			}
		}
		return nil, fmt.Errorf("token 校验失败")
	}, jwt.WithAudience(expectedAudience))

	if err != nil {
		return nil, fmt.Errorf("token 解析或校验失败: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("无效的 token")
	}

	return claims, nil
}