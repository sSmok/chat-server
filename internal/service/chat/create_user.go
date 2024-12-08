package chat

import (
	"context"

	"github.com/sSmok/chat-server/internal/model"
)

func (service *chatService) CreateUser(ctx context.Context, info *model.UserInfo) (int64, error) {
	return service.repo.CreateUser(ctx, info)
}
