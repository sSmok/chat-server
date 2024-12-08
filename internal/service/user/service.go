package user

import (
	"github.com/sSmok/chat-server/internal/repository"
	"github.com/sSmok/chat-server/internal/service"
)

type userService struct {
	repo repository.UserRepositoryI
}

// NewUserService создает объект сервиса для работы с пользователями на уровне сервисного слоя
func NewUserService(repo repository.UserRepositoryI) service.UserServiceI {
	return &userService{
		repo: repo,
	}
}
