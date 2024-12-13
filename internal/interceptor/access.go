package interceptor

import (
	"context"

	"github.com/sSmok/chat-server/internal/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var testAccessToken = "1"

type access struct {
	accessClient client.AccessCheckerI
}

// NewAccessInterceptor - конструктор интерцептора grpc-сервера для проверки прав доступа к указанному ендпоинту
func NewAccessInterceptor(accessClient client.AccessCheckerI) AccessInterceptorI {
	return &access{accessClient: accessClient}
}

func (a access) Access(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md := metadata.New(map[string]string{"Authorization": "Bearer " + testAccessToken})
	ctx = metadata.NewOutgoingContext(ctx, md)

	err := a.accessClient.Check(ctx, info.FullMethod)
	if err != nil {
		return nil, err
	}

	return handler(ctx, req)
}
