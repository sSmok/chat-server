package converter

import (
	"github.com/sSmok/chat-server/internal/model"
	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
)

// ToUserFromProto конвертирует данные пользователя из прото для сервисного слоя
func ToUserFromProto(user *descChat.User) *model.User {
	return &model.User{
		ID:   user.GetId(),
		Info: ToUserInfoFromProto(user.GetInfo()),
	}
}

// ToUserInfoFromProto конвертирует данные пользователя из прото для сервисного слоя
func ToUserInfoFromProto(info *descChat.UserInfo) model.UserInfo {
	return model.UserInfo{
		Name: info.GetName(),
	}
}

// ToProtoFromUser конвертирует данные пользователя из сервисного слоя в прото
func ToProtoFromUser(user *model.User) *descChat.User {
	return &descChat.User{
		Id:   user.ID,
		Info: ToProtoFromUserInfo(user.Info),
	}
}

// ToProtoFromUserInfo конвертирует данные пользователя из сервисного слоя в прото
func ToProtoFromUserInfo(info model.UserInfo) *descChat.UserInfo {
	return &descChat.UserInfo{
		Name: info.Name,
	}
}
