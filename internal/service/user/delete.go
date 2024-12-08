package user

import "context"

func (service *userService) DeleteUser(ctx context.Context, id int64) error {
	return service.repo.DeleteUser(ctx, id)
}
