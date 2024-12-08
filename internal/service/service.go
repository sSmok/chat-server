package service

import (
	"context"

	"github.com/sSmok/chat-server/internal/model"
)

// UserServiceI - интерфейс сервиса для работы с пользователями
type UserServiceI interface {
	CreateUser(ctx context.Context, info *model.UserInfo) (int64, error)
	DeleteUser(ctx context.Context, id int64) error
}

// ChatServiceI - интерфейс сервиса для работы с чатами
type ChatServiceI interface {
	CreateChat(ctx context.Context, chat *model.ChatInfo) (int64, error)
	DeleteChat(ctx context.Context, id int64) error
}

// MessageServiceI - интерфейс сервиса для работы с сообщениями
type MessageServiceI interface {
	CreateMessage(ctx context.Context, message *model.MessageInfo) error
}
