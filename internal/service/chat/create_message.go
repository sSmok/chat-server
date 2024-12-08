package chat

import (
	"context"

	"github.com/sSmok/chat-server/internal/model"
)

func (service *chatService) CreateMessage(ctx context.Context, message *model.MessageInfo) error {
	return service.repo.CreateMessage(ctx, message)
}
