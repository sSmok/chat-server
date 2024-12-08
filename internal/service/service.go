package service

import (
	"context"

	"github.com/sSmok/chat-server/internal/model"
)

// ChatServiceI - интерфейс сервиса для работы с чатами
type ChatServiceI interface {
	CreateChat(ctx context.Context, chat *model.ChatInfo) (int64, error)
	DeleteChat(ctx context.Context, id int64) error
	CreateUser(ctx context.Context, info *model.UserInfo) (int64, error)
	DeleteUser(ctx context.Context, id int64) error
	CreateMessage(ctx context.Context, message *model.MessageInfo) error
}
