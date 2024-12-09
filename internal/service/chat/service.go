package chat

import (
	"github.com/sSmok/chat-server/internal/client/db"
	"github.com/sSmok/chat-server/internal/repository"
	"github.com/sSmok/chat-server/internal/service"
)

type chatService struct {
	repo      repository.ChatRepositoryI
	txManager db.TxManagerI
}

// NewChatService создает объект сервиса для работы с чатами на уровне сервисного слоя
func NewChatService(repo repository.ChatRepositoryI, txManager db.TxManagerI) service.ChatServiceI {
	return &chatService{
		repo:      repo,
		txManager: txManager,
	}
}
