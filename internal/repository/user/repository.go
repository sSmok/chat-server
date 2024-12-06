package user

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sSmok/chat-server/internal/model"
	"github.com/sSmok/chat-server/internal/repository"
)

type userRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepo создает новый экземпляр репозитория пользователей
func NewUserRepo(pool *pgxpool.Pool) repository.UserRepositoryI {
	return &userRepo{pool: pool}
}

func (repo *userRepo) CreateUser(ctx context.Context, info *model.UserInfo) (int64, error) {
	userQuery := `insert into users (name) values (@name) returning id;`
	userQueryArgs := pgx.NamedArgs{
		"name": info.Name,
	}
	var userID int64
	err := repo.pool.QueryRow(ctx, userQuery, userQueryArgs).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (repo *userRepo) DeleteUser(ctx context.Context, id int64) error {
	userQuery := `delete from users where id=@id;`
	userQueryArgs := pgx.NamedArgs{
		"id": id,
	}
	_, err := repo.pool.Exec(ctx, userQuery, userQueryArgs)
	if err != nil {
		return err
	}

	return nil
}
