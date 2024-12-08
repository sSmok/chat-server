package message

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sSmok/chat-server/internal/model"
	"github.com/sSmok/chat-server/internal/repository"
)

type messageRepo struct {
	pool *pgxpool.Pool
}

// NewMessageRepo создает новый экземпляр репозитория сообщений
func NewMessageRepo(pool *pgxpool.Pool) repository.MessageRepositoryI {
	return &messageRepo{pool: pool}
}

func (repo *messageRepo) CreateMessage(ctx context.Context, info *model.MessageInfo) error {
	chatUserQuery := `select id from chats_users where chat_id=@chat_id and user_id=@user_id limit 1;`
	chatUserQueryArgs := pgx.NamedArgs{
		"chat_id": info.ChatID,
		"user_id": info.UserID,
	}
	var chatUserID int64
	err := repo.pool.QueryRow(ctx, chatUserQuery, chatUserQueryArgs).Scan(&chatUserID)
	if err != nil {
		return err
	}

	msgQuery := `insert into messages (source_id, text) values (@source_id, @text);`
	msgQueryArgs := pgx.NamedArgs{
		"source_id": chatUserID,
		"text":      info.Text,
	}

	_, err = repo.pool.Exec(ctx, msgQuery, msgQueryArgs)
	if err != nil {
		return err
	}

	return nil
}
