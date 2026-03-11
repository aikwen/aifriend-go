package auth

import (
	"context"

	"github.com/aikwen/aifriend-go/config"
	"github.com/aikwen/aifriend-go/internal/store"
	"github.com/aikwen/aifriend-go/internal/store/models"
)

type Service interface {
	Login(ctx context.Context, username, password string, jwtConf *config.JWTConfig) (*models.User, string, string, error)
	Register(ctx context.Context, username, password string, jwtConf *config.JWTConfig) (*models.User, string, string, error)
	RefreshToken(ctx context.Context, refreshTokenString string, jwtConf *config.JWTConfig) (string, string, error)
}


type authService struct {
	database *store.Database
}

func NewAuthService(database *store.Database) Service {
	return &authService{
		database: database,
	}
}
