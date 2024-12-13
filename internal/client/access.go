package client

import (
	"context"

	"github.com/sSmok/auth/pkg/access_v1"
)

type access struct {
	client access_v1.AccessV1Client
}

// NewAccessClient - создание клиента для проверки прав доступа
func NewAccessClient(client access_v1.AccessV1Client) AccessCheckerI {
	return &access{client: client}
}

func (a *access) Check(ctx context.Context, endpoint string) error {
	_, err := a.client.Check(ctx, &access_v1.CheckRequest{EndpointAddress: endpoint})
	if err != nil {
		return err
	}

	return nil
}
