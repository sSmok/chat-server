package chat

import (
	"context"

	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteChat удаляет чат по его ID
func (api *API) DeleteChat(ctx context.Context, req *descChat.DeleteChatRequest) (*emptypb.Empty, error) {
	err := api.chatService.DeleteChat(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
