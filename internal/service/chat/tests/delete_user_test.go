package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/sSmok/chat-server/internal/repository"
	"github.com/sSmok/chat-server/internal/service/chat"
	"github.com/sSmok/platform_common/pkg/client/db/transaction"
	"github.com/stretchr/testify/require"

	"github.com/sSmok/chat-server/internal/repository/mocks"
	"github.com/sSmok/platform_common/pkg/client/db"
	txMocks "github.com/sSmok/platform_common/pkg/client/db/mocks"
)

func Test_chatService_DeleteUser(t *testing.T) {
	type chatRepositoryIMockFunc func(mc *minimock.Controller) repository.ChatRepositoryI
	type transactorIMockFunc func(mc *minimock.Controller) db.TransactorI

	type args struct {
		ctx context.Context
		id  int64
	}

	var (
		ctx           = context.Background()
		minimockContr = minimock.NewController(t)
		id            = gofakeit.Int64()
		repoErr       = errors.New("error")
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
				ctx: ctx,
				id:  id,
			},
			want: nil,
			chatRepositoryIMock: func(mc *minimock.Controller) repository.ChatRepositoryI {
				mock := mocks.NewChatRepositoryIMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(nil)
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
				ctx: ctx,
				id:  id,
			},
			want: repoErr,
			chatRepositoryIMock: func(mc *minimock.Controller) repository.ChatRepositoryI {
				mock := mocks.NewChatRepositoryIMock(mc)
				mock.DeleteUserMock.Expect(ctx, id).Return(repoErr)
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

			err := chatServ.DeleteUser(tt.args.ctx, tt.args.id)
			require.Equal(t, tt.want, err)
		})
	}
}
