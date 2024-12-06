package converter

import (
	"github.com/sSmok/chat-server/internal/model"
	modelRepo "github.com/sSmok/chat-server/internal/repository/message/model"
)

// ToRepoFromMessage преобразует сообщение из сервисной модели в модель репозитория
func ToRepoFromMessage(message *model.Message) *modelRepo.Message {
	return &modelRepo.Message{
		ID:   message.ID,
		Info: ToRepoFromMessageInfo(message.Info),
	}
}

// ToRepoFromMessageInfo преобразует информацию о сообщении из сервисной модели в модель репозитория
func ToRepoFromMessageInfo(messageInfo model.MessageInfo) modelRepo.MessageInfo {
	return modelRepo.MessageInfo{
		ChatID: messageInfo.ChatID,
		UserID: messageInfo.UserID,
		Text:   messageInfo.Text,
	}
}

// ToMessageInfoFromRepo преобразует информацию о сообщении из репозитория в сервисную модель
func ToMessageInfoFromRepo(info modelRepo.MessageInfo) model.MessageInfo {
	return model.MessageInfo{
		ChatID: info.ChatID,
		UserID: info.UserID,
		Text:   info.Text,
	}
}
