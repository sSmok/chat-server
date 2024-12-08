package message

import (
	"github.com/sSmok/chat-server/internal/repository"
	"github.com/sSmok/chat-server/internal/service"
)

type messageService struct {
	repo repository.MessageRepositoryI
}

// NewMessageService создает новый сервис для работы с сообщениями
func NewMessageService(repo repository.MessageRepositoryI) service.MessageServiceI {
	return &messageService{
		repo: repo,
	}
}
