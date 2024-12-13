package client

import "context"

// AccessCheckerI - интерфейс проверки прав доступа к указанному ендпоинту
type AccessCheckerI interface {
	Check(ctx context.Context, endpoint string) error
}
