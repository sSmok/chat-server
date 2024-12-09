package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/sSmok/chat-server/internal/api/chat"
	"github.com/sSmok/chat-server/internal/model"
	"github.com/sSmok/chat-server/internal/service"
	"github.com/sSmok/chat-server/internal/service/mocks"
	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestAPI_CreateMessage(t *testing.T) {
	type userServiceIMockFunc func(mc *minimock.Controller) service.ChatServiceI

	type args struct {
		ctx context.Context
		req *descChat.CreateMessageRequest
	}

	var (
		ctx           = context.Background()
		minimockContr = minimock.NewController(t)
		userID        = gofakeit.Int64()
		chatID        = gofakeit.Int64()
		text          = gofakeit.Sentence(5)
		serviceErr    = fmt.Errorf("service error")

		req = &descChat.CreateMessageRequest{
			Info: &descChat.MessageInfo{
				UserId: userID,
				ChatId: chatID,
				Text:   text,
			},
		}

		info = &model.MessageInfo{
			ChatID: chatID,
			UserID: userID,
			Text:   text,
		}
	)

	tests := []struct {
		name             string
		args             args
		want             *emptypb.Empty
		err              error
		userServiceIMock userServiceIMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: &emptypb.Empty{},
			err:  nil,
			userServiceIMock: func(mc *minimock.Controller) service.ChatServiceI {
				mock := mocks.NewChatServiceIMock(mc)
				mock.CreateMessageMock.Expect(ctx, info).Return(nil)
				return mock
			},
		},
		{
			name: "fail case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceErr,
			userServiceIMock: func(mc *minimock.Controller) service.ChatServiceI {
				mock := mocks.NewChatServiceIMock(mc)
				mock.CreateMessageMock.Expect(ctx, info).Return(serviceErr)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceIMock := tt.userServiceIMock(minimockContr)
			api := chat.NewAPI(userServiceIMock)

			newResp, err := api.CreateMessage(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newResp)
		})
	}
}
