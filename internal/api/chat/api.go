package chat

import (
	"github.com/sSmok/chat-server/internal/service"
	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
)

// API - апи слой для работы с чатом, взаимодействует с сервисным слоем
type API struct {
	descChat.UnimplementedChatV1Server
	chatService service.ChatServiceI
}

// NewAPI - конструктор апи слоя
func NewAPI(chatService service.ChatServiceI) *API {
	return &API{
		chatService: chatService,
	}
}
