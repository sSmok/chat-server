package converter

import (
	"github.com/sSmok/chat-server/internal/model"
	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
)

// ToChatFromProto конвертирует данные чата из прото для сервисного слоя
func ToChatFromProto(chat *descChat.Chat) *model.Chat {
	return &model.Chat{
		ID:   chat.GetId(),
		Info: ToChatInfoFromProto(chat.GetInfo()),
	}
}

// ToChatInfoFromProto конвертирует данные чата из прото для сервисного слоя
func ToChatInfoFromProto(info *descChat.ChatInfo) model.ChatInfo {
	return model.ChatInfo{
		Name:    info.GetName(),
		UserIDs: info.UserIds,
	}
}

// ToProtoFromChat конвертирует данные чата из сервисного слоя в прото
func ToProtoFromChat(chat *model.Chat) *descChat.Chat {
	return &descChat.Chat{
		Id:   chat.ID,
		Info: ToProtoFromChatInfo(chat.Info),
	}
}

// ToProtoFromChatInfo конвертирует данные чата из сервисного слоя в прото
func ToProtoFromChatInfo(info model.ChatInfo) *descChat.ChatInfo {
	return &descChat.ChatInfo{
		Name: info.Name,
	}
}
