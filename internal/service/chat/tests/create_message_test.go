package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/sSmok/chat-server/internal/model"
	"github.com/sSmok/chat-server/internal/repository"
	"github.com/sSmok/chat-server/internal/service/chat"
	"github.com/sSmok/platform_common/pkg/client/db/transaction"
	"github.com/stretchr/testify/require"

	"github.com/sSmok/chat-server/internal/repository/mocks"
	"github.com/sSmok/platform_common/pkg/client/db"
	txMocks "github.com/sSmok/platform_common/pkg/client/db/mocks"
)

func Test_chatService_CreateMessage(t *testing.T) {
	type chatRepositoryIMockFunc func(mc *minimock.Controller) repository.ChatRepositoryI
	type transactorIMockFunc func(mc *minimock.Controller) db.TransactorI

	type args struct {
		ctx  context.Context
		info *model.MessageInfo
	}

	var (
		ctx           = context.Background()
		minimockContr = minimock.NewController(t)
		chatID        = gofakeit.Int64()
		userID        = gofakeit.Int64()
		text          = gofakeit.Sentence(5)
		repoErr       = errors.New("error")

		info = &model.MessageInfo{
			ChatID: chatID,
			UserID: userID,
			Text:   text,
		}
	)

	tests := []struct {
		name                string
		args                args
		want                error
		chatRepositoryIMock chatRepositoryIMockFunc
		transactorIMock     transactorIMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:  ctx,
				info: info,
			},
			want: nil,
			chatRepositoryIMock: func(mc *minimock.Controller) repository.ChatRepositoryI {
				mock := mocks.NewChatRepositoryIMock(mc)
				mock.CreateMessageMock.Expect(ctx, info).Return(nil)
				return mock
			},
			transactorIMock: func(mc *minimock.Controller) db.TransactorI {
				mock := txMocks.NewTransactorIMock(mc)
				return mock
			},
		},
		{
			name: "fail case",
			args: args{
				ctx:  ctx,
				info: info,
			},
			want: repoErr,
			chatRepositoryIMock: func(mc *minimock.Controller) repository.ChatRepositoryI {
				mock := mocks.NewChatRepositoryIMock(mc)
				mock.CreateMessageMock.Expect(ctx, info).Return(repoErr)
				return mock
			},
			transactorIMock: func(mc *minimock.Controller) db.TransactorI {
				mock := txMocks.NewTransactorIMock(mc)
				return mock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chatRepositoryIMock := tt.chatRepositoryIMock(minimockContr)
			txManager := transaction.NewManager(tt.transactorIMock(minimockContr))
			chatServ := chat.NewChatService(chatRepositoryIMock, txManager)

			err := chatServ.CreateMessage(tt.args.ctx, tt.args.info)
			require.Equal(t, tt.want, err)
		})
	}
}
