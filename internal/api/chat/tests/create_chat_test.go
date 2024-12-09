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
)

func TestAPI_CreateChat(t *testing.T) {
	type userServiceIMockFunc func(mc *minimock.Controller) service.ChatServiceI

	type args struct {
		ctx context.Context
		req *descChat.CreateChatRequest
	}

	var (
		ctx           = context.Background()
		minimockContr = minimock.NewController(t)
		id            = gofakeit.Int64()
		name          = gofakeit.Name()
		userIDs       = []int64{gofakeit.Int64(), gofakeit.Int64()}
		serviceErr    = fmt.Errorf("service error")

		req = &descChat.CreateChatRequest{
			Info: &descChat.ChatInfo{
				Name:    name,
				UserIds: userIDs,
			},
		}

		info = &model.ChatInfo{
			Name:    name,
			UserIDs: userIDs,
		}

		resp = &descChat.CreateChatResponse{Id: id}
	)

	tests := []struct {
		name             string
		args             args
		want             *descChat.CreateChatResponse
		err              error
		userServiceIMock userServiceIMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: resp,
			err:  nil,
			userServiceIMock: func(mc *minimock.Controller) service.ChatServiceI {
				mock := mocks.NewChatServiceIMock(mc)
				mock.CreateChatMock.Expect(ctx, info).Return(id, nil)
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
				mock.CreateChatMock.Expect(ctx, info).Return(0, serviceErr)
				return mock
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			userServiceIMock := tt.userServiceIMock(minimockContr)
			api := chat.NewAPI(userServiceIMock)

			newResp, err := api.CreateChat(tt.args.ctx, tt.args.req)
			require.Equal(t, tt.err, err)
			require.Equal(t, tt.want, newResp)
		})
	}
}
