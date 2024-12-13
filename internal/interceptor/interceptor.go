package interceptor

import (
	"context"

	"google.golang.org/grpc"
)

// AccessInterceptorI - интерфейс интерцептора grpc-сервера для проверки прав доступа к указанному ендпоинту
type AccessInterceptorI interface {
	Access(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}
