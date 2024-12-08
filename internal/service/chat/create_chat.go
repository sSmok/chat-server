package chat

import (
	"context"

	"github.com/sSmok/chat-server/internal/model"
)

func (service *chatService) CreateChat(ctx context.Context, chat *model.ChatInfo) (int64, error) {
	return service.repo.CreateChat(ctx, chat)
}
