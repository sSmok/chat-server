package chat

import (
	"github.com/sSmok/chat-server/internal/service"
	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
)

// API - апи слой для работы с чатом, взаимодействует с сервисным слоем
type API struct {
	descChat.UnimplementedChatV1Server
	chatService    service.ChatServiceI
	userService    service.UserServiceI
	messageService service.MessageServiceI
}

// NewAPI - конструктор апи слоя
func NewAPI(chatService service.ChatServiceI, userService service.UserServiceI, messageService service.MessageServiceI) *API {
	return &API{
		chatService:    chatService,
		userService:    userService,
		messageService: messageService,
	}
}
