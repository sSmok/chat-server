package converter

import (
	"github.com/sSmok/chat-server/internal/model"
	modelRepo "github.com/sSmok/chat-server/internal/repository/chat/model"
)

// ToUserFromRepo преобразует модель пользователя из репозитория в сервисную модель
func ToUserFromRepo(userRepo *modelRepo.User) *model.User {
	return &model.User{
		ID:   userRepo.ID,
		Info: ToUserInfoFromRepo(userRepo.Info),
	}
}

// ToUserInfoFromRepo преобразует информацию о пользователе из репозитория в сервисную модель
func ToUserInfoFromRepo(userInfoRepo modelRepo.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name: userInfoRepo.Name,
	}
}

// ToRepoFromUserInfo преобразует информацию о пользователе из сервисной модели в модель репозитория
func ToRepoFromUserInfo(user *model.UserInfo) *modelRepo.UserInfo {
	return &modelRepo.UserInfo{
		Name: user.Name,
	}
}
