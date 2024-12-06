package chat

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sSmok/chat-server/internal/repository"
)

type chatRepo struct {
	pool *pgxpool.Pool
}

// NewChatRepo создает новый экземпляр репозитория чатов
func NewChatRepo(pool *pgxpool.Pool) repository.ChatRepositoryI {
	return &chatRepo{pool: pool}
}

func (repo *chatRepo) CreateChat(ctx context.Context, name string, userIDs []int64) (int64, error) {
	var chatID int64
	tx, err := repo.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	chatQuery := `insert into chats (name) values (@name) returning id;`
	chatQueryArgs := pgx.NamedArgs{
		"name": name,
	}

	row := tx.QueryRow(ctx, chatQuery, chatQueryArgs)
	err = row.Scan(&chatID)
	if err != nil {
		return 0, err
	}

	batch := &pgx.Batch{}
	chatUserQuery := `insert into chats_users (chat_id, user_id) values (@chat_id, @user_id);`
	for _, id := range userIDs {
		chatUserQueryArgs := pgx.NamedArgs{
			"chat_id": chatID,
			"user_id": id,
		}
		batch.Queue(chatUserQuery, chatUserQueryArgs)
	}

	r := tx.SendBatch(ctx, batch)
	err = r.Close()
	if err != nil {
		return 0, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

func (repo *chatRepo) DeleteChat(ctx context.Context, id int64) error {
	delChat := `delete from chats where id = @id`
	delChatArgs := pgx.NamedArgs{
		"id": id,
	}

	_, err := repo.pool.Exec(ctx, delChat, delChatArgs)
	if err != nil {
		return err
	}

	return nil
}
