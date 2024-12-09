package chat

import (
	"context"

	"github.com/sSmok/chat-server/internal/model"
)

func (service *chatService) CreateChat(ctx context.Context, chat *model.ChatInfo) (int64, error) {
	var chatID int64
	err := service.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		chatID, errTx = service.repo.CreateChat(ctx, chat)
		if errTx != nil {
			return errTx
		}

		errTx = service.repo.AddUsersToChat(ctx, chatID, chat.UserIDs)
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return chatID, nil
}
