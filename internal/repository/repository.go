package repository

import (
	"context"

	"github.com/sSmok/chat-server/internal/model"
)

// ChatRepositoryI интерфейс предоставляет методы для работы с чатами
type ChatRepositoryI interface {
	CreateChat(ctx context.Context, info *model.ChatInfo) (int64, error)
	DeleteChat(ctx context.Context, id int64) error
	CreateUser(ctx context.Context, info *model.UserInfo) (int64, error)
	DeleteUser(ctx context.Context, id int64) error
	CreateMessage(ctx context.Context, info *model.MessageInfo) error
}
