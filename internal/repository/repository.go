package repository

import (
	"context"

	"github.com/sSmok/chat-server/internal/model"
)

// UserRepositoryI интерфейс предоставляет методы для работы с пользователями
type UserRepositoryI interface {
	CreateUser(ctx context.Context, info *model.UserInfo) (int64, error)
	DeleteUser(ctx context.Context, id int64) error
}

// ChatRepositoryI интерфейс предоставляет методы для работы с чатами
type ChatRepositoryI interface {
	CreateChat(ctx context.Context, name string, userIDs []int64) (int64, error)
	DeleteChat(ctx context.Context, id int64) error
}

// MessageRepositoryI интерфейс предоставляет методы для работы с сообщениями
type MessageRepositoryI interface {
	CreateMessage(ctx context.Context, chatID, userID int64, text string) error
}
