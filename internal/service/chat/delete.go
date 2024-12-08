package chat

import "context"

func (service *chatService) DeleteChat(ctx context.Context, id int64) error {
	return service.repo.DeleteChat(ctx, id)
}
