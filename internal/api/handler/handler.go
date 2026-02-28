package handler

import (
	"github.com/aikwen/aifriend-go/internal/auth"
	"github.com/aikwen/aifriend-go/internal/character"
	"github.com/aikwen/aifriend-go/internal/user"
	"github.com/aikwen/aifriend-go/pkg/storage"
)

type Handler struct {
	authSvc auth.Service
	charSvc character.Service
	userSvc user.Service
	storage storage.Storage
}

func NewHandler(authSvc auth.Service,
	charSvc character.Service,
	userSvc user.Service,
	st storage.Storage) *Handler {
	return &Handler{
		authSvc: authSvc,
		charSvc: charSvc,
		userSvc: userSvc,
		storage: st,
	}
}
