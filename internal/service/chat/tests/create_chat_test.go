package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/jackc/pgx/v5"
	"github.com/sSmok/chat-server/internal/model"
	"github.com/sSmok/chat-server/internal/repository"
	"github.com/sSmok/chat-server/internal/service/chat"
	"github.com/sSmok/platform_common/pkg/client/db/pg"
	"github.com/sSmok/platform_common/pkg/client/db/transaction"
	"github.com/stretchr/testify/require"

	"github.com/sSmok/chat-server/internal/repository/mocks"
	"github.com/sSmok/platform_common/pkg/client/db"
	txMocks "github.com/sSmok/platform_common/pkg/client/db/mocks"
)

func Test_chatService_CreateChat(t *testing.T) {
	type chatRepositoryIMockFunc func(mc *minimock.Controller) repository.ChatRepositoryI
	type transactorIMockFunc func(mc *minimock.Controller) db.TransactorI

	type args struct {
		ctx  context.Context
		info *model.ChatInfo
	}

	var (
		ctx           = context.Background()
		minimockContr = minimock.NewController(t)
		id            = gofakeit.Int64()
		userIDs       = []int64{gofakeit.Int64(), gofakeit.Int64()}
		name          = gofakeit.Name()
		repoErr       = errors.New("error")

		info = &model.ChatInfo{
			Name:    name,
			UserIDs: userIDs,
		}

		txMock txMocks.TxMock
	)

	tests := []struct {
		name                string
		args                args
		want                int64
		err                 error
		chatRepositoryIMock chatRepositoryIMockFunc
		transactorIMock     transactorIMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:  ctx,
				info: info,
			},
			want: id,
			err:  nil,
			chatRepositoryIMock: func(mc *minimock.Controller) repository.ChatRepositoryI {
				ctxTx := pg.MakeContextTransaction(ctx, &txMock)
				mock := mocks.NewChatRepositoryIMock(mc)
				mock.CreateChatMock.Expect(ctxTx, info).Return(id, nil)
				mock.AddUsersToChatMock.Expect(ctxTx, id, userIDs).Return(nil)
				return mock
			},
			transactorIMock: func(mc *minimock.Controller) db.TransactorI {
				mock := txMocks.NewTransactorIMock(mc)
				mock.BeginTxMock.Expect(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).Return(&txMock, nil)
				return mock
			},
		},
		{
			name: "fail case",
			args: args{
				ctx:  ctx,
				info: info,
			},
			want: 0,
			err:  repoErr,
			chatRepositoryIMock: func(mc *minimock.Controller) repository.ChatRepositoryI {
				ctxTx := pg.MakeContextTransaction(ctx, &txMock)
				mock := mocks.NewChatRepositoryIMock(mc)
				mock.CreateChatMock.Expect(ctxTx, info).Return(0, repoErr)
				//mock.AddUsersToChatMock.Expect(ctxTx, id, userIDs).Return(nil)
				return mock
			},
			transactorIMock: func(mc *minimock.Controller) db.TransactorI {
				mock := txMocks.NewTransactorIMock(mc)
				mock.BeginTxMock.Expect(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted}).Return(&txMock, nil)
				return mock
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chatRepositoryIMock := tt.chatRepositoryIMock(minimockContr)
			txManager := transaction.NewManager(tt.transactorIMock(minimockContr))
			chatServ := chat.NewChatService(chatRepositoryIMock, txManager)

			userID, err := chatServ.CreateChat(tt.args.ctx, tt.args.info)

			if tt.name == "success case" {
				require.Equal(t, tt.err, err)
			} else {
				require.Contains(t, err.Error(), tt.err.Error())
			}

			require.Equal(t, tt.want, userID)
		})
	}
}
