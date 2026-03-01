package friend

import (
	"path"

	"github.com/aikwen/aifriend-go/internal/models"
	"gorm.io/gorm"
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
	GetOrCreate(userID uint, characterID uint) (*FriendDTO, error)
	GetList(userID uint, offset int, limit int) ([]FriendDTO, error)
	Remove(userID uint, friendID uint) error
}


type friendService struct {
	store store
}

// NewService 构造函数
func NewService(db *gorm.DB) Service {
	return &friendService{
		store: newFriendStore(db),
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