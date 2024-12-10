package chat

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/sSmok/chat-server/internal/model"
	"github.com/sSmok/chat-server/internal/repository"
	"github.com/sSmok/platform_common/pkg/client/db"
	"github.com/sSmok/platform_common/pkg/client/db/pg"
)

type chatRepo struct {
	dbClient db.ClientI
}

// NewChatRepo создает новый экземпляр репозитория чатов
func NewChatRepo(dbClient db.ClientI) repository.ChatRepositoryI {
	return &chatRepo{dbClient: dbClient}
}

func (repo *chatRepo) CreateChat(ctx context.Context, info *model.ChatInfo) (int64, error) {
	var chatID int64
	chatQuery := `insert into chats (name) values (@name) returning id;`
	chatQueryArgs := pgx.NamedArgs{
		"name": info.Name,
	}
	q := db.Query{
		Name:     "chat_repository.CreateChat",
		QueryRaw: chatQuery,
	}

	err := repo.dbClient.DB().ScanOneContext(ctx, &chatID, q, chatQueryArgs)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}

func (repo *chatRepo) AddUsersToChat(ctx context.Context, chatID int64, userIDs []int64) error {
	batch := &pgx.Batch{}
	chatUserQuery := `insert into chats_users (chat_id, user_id) values (@chat_id, @user_id);`
	for _, id := range userIDs {
		chatUserQueryArgs := pgx.NamedArgs{
			"chat_id": chatID,
			"user_id": id,
		}
		batch.Queue(chatUserQuery, chatUserQueryArgs)
	}

	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if !ok {
		return errors.New("cant make batch request, transaction not present")
	}

	r := tx.SendBatch(ctx, batch)
	err := r.Close()
	if err != nil {
		return err
	}

	return nil
}

func (repo *chatRepo) DeleteChat(ctx context.Context, id int64) error {
	query := `delete from chats where id = @id`
	queryArgs := pgx.NamedArgs{
		"id": id,
	}

	q := db.Query{
		Name:     "chat_repository.DeleteChat",
		QueryRaw: query,
	}

	_, err := repo.dbClient.DB().ExecContext(ctx, q, queryArgs)
	if err != nil {
		return err
	}

	return nil
}

func (repo *chatRepo) CreateUser(ctx context.Context, info *model.UserInfo) (int64, error) {
	query := `insert into users (name) values (@name) returning id;`
	queryArgs := pgx.NamedArgs{
		"name": info.Name,
	}
	q := db.Query{
		Name:     "user_repository.CreateUser",
		QueryRaw: query,
	}

	var userID int64
	err := repo.dbClient.DB().ScanOneContext(ctx, &userID, q, queryArgs)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (repo *chatRepo) DeleteUser(ctx context.Context, id int64) error {
	query := `delete from users where id=@id;`
	queryArgs := pgx.NamedArgs{
		"id": id,
	}
	q := db.Query{
		Name:     "user_repository.DeleteUser",
		QueryRaw: query,
	}

	_, err := repo.dbClient.DB().ExecContext(ctx, q, queryArgs)
	if err != nil {
		return err
	}

	return nil
}

func (repo *chatRepo) CreateMessage(ctx context.Context, info *model.MessageInfo) error {
	queryChatUser := `select id from chats_users where chat_id=@chat_id and user_id=@user_id limit 1;`
	queryChatUserArgs := pgx.NamedArgs{
		"chat_id": info.ChatID,
		"user_id": info.UserID,
	}
	qChatUser := db.Query{
		Name:     "chat_repository.CreateMessage.SelectChatUser",
		QueryRaw: queryChatUser,
	}

	var chatUserID int64
	err := repo.dbClient.DB().ScanOneContext(ctx, &chatUserID, qChatUser, queryChatUserArgs)
	if err != nil {
		return err
	}

	msgQuery := `insert into messages (source_id, text) values (@source_id, @text);`
	msgQueryArgs := pgx.NamedArgs{
		"source_id": chatUserID,
		"text":      info.Text,
	}
	qMsg := db.Query{
		Name:     "chat_repository.CreateMessage.InsertMessage",
		QueryRaw: msgQuery,
	}

	_, err = repo.dbClient.DB().ExecContext(ctx, qMsg, msgQueryArgs)
	if err != nil {
		return err
	}

	return nil
}
