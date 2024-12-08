package message

import (
	"context"

	"github.com/sSmok/chat-server/internal/model"
)

func (service *messageService) CreateMessage(ctx context.Context, message *model.MessageInfo) error {
	return service.repo.CreateMessage(ctx, message)
}
