package chat

import (
	"context"

	"github.com/sSmok/chat-server/internal/converter"
	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
)

// CreateChat - создание чата
func (api *API) CreateChat(ctx context.Context, req *descChat.CreateChatRequest) (*descChat.CreateChatResponse, error) {
	info := converter.ToChatInfoFromProto(req.GetInfo())
	chatID, err := api.chatService.CreateChat(ctx, &info)
	if err != nil {
		return nil, err
	}

	return &descChat.CreateChatResponse{Id: chatID}, nil
}
