package auth

import (
	"context"

	"github.com/aikwen/aifriend-go/config"
	"github.com/aikwen/aifriend-go/internal/models"
)

type Service interface {
	Login(ctx context.Context, username, password string) (*models.User, string, string, error)
	Register(ctx context.Context, username, password string) (*models.User, string, string, error)
	RefreshToken(ctx context.Context, refreshTokenString string) (string, string, error)
}

type userService interface {
	Create(ctx context.Context, user *models.User) error
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
}

type authService struct {
	userService userService
	jwtConf     *config.JWTConfig
}

func NewAuthService(us userService, jwtConf *config.JWTConfig) Service {
	return &authService{
		userService: us,
		jwtConf:     jwtConf,
	}
}
