package chat

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sSmok/chat-server/internal/model"
	"github.com/sSmok/chat-server/internal/repository"
)

type chatRepo struct {
	pool *pgxpool.Pool
}

// NewChatRepo создает новый экземпляр репозитория чатов
func NewChatRepo(pool *pgxpool.Pool) repository.ChatRepositoryI {
	return &chatRepo{pool: pool}
}

func (repo *chatRepo) CreateChat(ctx context.Context, info *model.ChatInfo) (int64, error) {
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
		"name": info.Name,
	}

	row := tx.QueryRow(ctx, chatQuery, chatQueryArgs)
	err = row.Scan(&chatID)
	if err != nil {
		return 0, err
	}

	batch := &pgx.Batch{}
	chatUserQuery := `insert into chats_users (chat_id, user_id) values (@chat_id, @user_id);`
	for _, id := range info.UserIDs {
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

func (repo *chatRepo) CreateUser(ctx context.Context, info *model.UserInfo) (int64, error) {
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

func (repo *chatRepo) DeleteUser(ctx context.Context, id int64) error {
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

func (repo *chatRepo) CreateMessage(ctx context.Context, info *model.MessageInfo) error {
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
