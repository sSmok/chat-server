package converter

import (
	"github.com/sSmok/chat-server/internal/model"
	modelRepo "github.com/sSmok/chat-server/internal/repository/chat/model"
)

// ToChatFromRepo преобразует модель чата из репозитория в сервисную модель
func ToChatFromRepo(chat *modelRepo.Chat) *model.Chat {
	return &model.Chat{
		ID:   chat.ID,
		Info: ToChatInfoFromRepo(chat.Info),
	}
}

// ToChatInfoFromRepo преобразует информацию о чате из репозитория в сервисную модель
func ToChatInfoFromRepo(info modelRepo.ChatInfo) model.ChatInfo {
	return model.ChatInfo{
		Name:    info.Name,
		UserIDs: info.UserIDs,
	}
}

// ToRepoFromChatInfo преобразует модель данных чата из сервисной модели в модель репозитория
func ToRepoFromChatInfo(chatInfo *model.ChatInfo) *modelRepo.ChatInfo {
	return &modelRepo.ChatInfo{
		Name:    chatInfo.Name,
		UserIDs: chatInfo.UserIDs,
	}
}
