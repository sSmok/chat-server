package chat

import (
	"context"

	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteUser удаляет пользователя по его ID
func (api *API) DeleteUser(ctx context.Context, req *descChat.DeleteUserRequest) (*emptypb.Empty, error) {
	err := api.chatService.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
