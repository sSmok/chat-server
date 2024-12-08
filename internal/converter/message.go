package converter

import (
	"github.com/sSmok/chat-server/internal/model"
	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
)

// ToMessageFromProto конвертирует данные сообщения из прото для сервисного слоя
func ToMessageFromProto(message *descChat.Message) *model.Message {
	return &model.Message{
		ID:   message.GetId(),
		Info: ToMessageInfoFromProto(message.GetInfo()),
	}
}

// ToMessageInfoFromProto конвертирует данные сообщения из прото для сервисного слоя
func ToMessageInfoFromProto(info *descChat.MessageInfo) model.MessageInfo {
	return model.MessageInfo{
		ChatID: info.GetChatId(),
		UserID: info.GetUserId(),
		Text:   info.GetText(),
	}
}
