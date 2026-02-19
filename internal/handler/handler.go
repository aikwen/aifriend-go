package handler


import (
	"github.com/aikwen/aifriend-go/internal/service"
)

// Handler 是所有具体业务 Handler 的总聚合体
type Handler struct {
	Auth *AuthHandler
	User *UserHandler
}

// NewHandler 负责接收所有的 Service，并统一初始化所有的 Handler
func NewHandler(authSvc service.AuthService, userSvc service.UserService) *Handler {
	return &Handler{
		Auth: NewAuthHandler(authSvc),
		User: NewUserHandler(userSvc),
	}
}