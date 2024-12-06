package message

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sSmok/chat-server/internal/repository"
)

type messageRepo struct {
	pool *pgxpool.Pool
}

// NewMessageRepo создает новый экземпляр репозитория сообщений
func NewMessageRepo(pool *pgxpool.Pool) repository.MessageRepositoryI {
	return &messageRepo{pool: pool}
}

func (repo *messageRepo) CreateMessage(ctx context.Context, chatID, userID int64, text string) error {
	chatUserQuery := `select id from chats_users where chat_id=@chat_id and user_id=@user_id limit 1;`
	chatUserQueryArgs := pgx.NamedArgs{
		"chat_id": chatID,
		"user_id": userID,
	}
	var chatUserID int64
	err := repo.pool.QueryRow(ctx, chatUserQuery, chatUserQueryArgs).Scan(&chatUserID)
	if err != nil {
		return err
	}

	msgQuery := `insert into messages (source_id, text) values (@source_id, @text) returning id;`
	msgQueryArgs := pgx.NamedArgs{
		"source_id": chatUserID,
		"text":      text,
	}

	_, err = repo.pool.Exec(ctx, msgQuery, msgQueryArgs)
	if err != nil {
		return err
	}

	return nil
}
