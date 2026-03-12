package friend

import (
	"context"
	"path"

	"github.com/aikwen/aifriend-go/internal/store"
	"github.com/aikwen/aifriend-go/internal/store/models"
)

type AuthorDTO struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Photo    string `json:"photo"`
}

type CharacterDTO struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Profile         string    `json:"profile"`
	Photo           string    `json:"photo"`
	BackgroundImage string    `json:"background_image"`
	Author          AuthorDTO `json:"author"`
}

type FriendDTO struct {
	ID        uint         `json:"id"`
	Character CharacterDTO `json:"character"`
}

type Service interface {
	GetOrCreate(ctx context.Context, userID uint, characterID uint) (*FriendDTO, error)
	GetList(ctx context.Context, userID uint, offset int, limit int) ([]FriendDTO, error)
	Remove(ctx context.Context, userID uint, friendID uint) error
}

type friendService struct {
	database *store.Database
}

// NewFriendService 构造函数
func NewFriendService(database *store.Database) Service {
	return &friendService{
		database: database,
	}
}

func (s *friendService) convertToDTO(friend *models.Friend) FriendDTO {
	return FriendDTO{
		ID: friend.ID,
		Character: CharacterDTO{
			ID:              friend.Character.ID,
			Name:            friend.Character.Name,
			Profile:         friend.Character.Profile,
			Photo:           path.Join("/media", friend.Character.Photo),
			BackgroundImage: path.Join("/media", friend.Character.BackgroundImage),
			Author: AuthorDTO{
				UserID:   friend.Character.Author.ID,
				Username: friend.Character.Author.Username,
				Photo:    path.Join("/media", friend.Character.Author.Photo),
			},
		},
	}
}
