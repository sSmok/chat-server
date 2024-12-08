package chat

import "context"

func (service *chatService) DeleteUser(ctx context.Context, id int64) error {
	return service.repo.DeleteUser(ctx, id)
}
