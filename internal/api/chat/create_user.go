package chat

import (
	"context"

	"github.com/sSmok/chat-server/internal/converter"
	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
)

// CreateUser - создание пользователя
func (api *API) CreateUser(ctx context.Context, req *descChat.CreateUserRequest) (*descChat.CreateUserResponse, error) {
	info := converter.ToUserInfoFromProto(req.GetInfo())
	userID, err := api.chatService.CreateUser(ctx, &info)
	if err != nil {
		return nil, err
	}
	return &descChat.CreateUserResponse{Id: userID}, nil
}
