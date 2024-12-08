package chat

import (
	"context"

	"github.com/sSmok/chat-server/internal/converter"
	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateMessage - создание сообщения
func (api *API) CreateMessage(ctx context.Context, req *descChat.CreateMessageRequest) (*emptypb.Empty, error) {
	info := converter.ToMessageInfoFromProto(req.GetInfo())
	err := api.messageService.CreateMessage(ctx, &info)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
