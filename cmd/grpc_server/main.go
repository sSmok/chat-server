package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sSmok/chat-server/internal/config"
	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config file: %v", err)
	}
	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to load GRPC config")
	}
	pgConfig, err := config.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to load GRPC config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %+v\n", err)
	}
	defer func() {
		if err := lis.Close(); err != nil {
			log.Fatalf("listener cannot be closed: %v", err)
		}
	}()

	ctx := context.Background()
	serv := grpc.NewServer()
	reflection.Register(serv)
	pool, err := pgxpool.New(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to load PG config: %v", err)
	}

	descChat.RegisterChatV1Server(serv, &server{pool: pool})
	if err = serv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %+v\n", err)
	}
}

type server struct {
	pool *pgxpool.Pool
	descChat.UnimplementedChatV1Server
}

func (s *server) CreateUser(ctx context.Context, req *descChat.CreateUserRequest) (*descChat.CreateUserResponse, error) {
	userQuery := `insert into users (name) values (@name) returning id;`
	userQueryArgs := pgx.NamedArgs{
		"name": req.GetName(),
	}
	var userID int64
	err := s.pool.QueryRow(ctx, userQuery, userQueryArgs).Scan(&userID)
	if err != nil {
		return nil, err
	}

	return &descChat.CreateUserResponse{Id: userID}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *descChat.DeleteUserRequest) (*emptypb.Empty, error) {
	userQuery := `delete from users where id=@id;`
	userQueryArgs := pgx.NamedArgs{
		"id": req.GetId(),
	}
	_, err := s.pool.Exec(ctx, userQuery, userQueryArgs)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) CreateChat(ctx context.Context, req *descChat.CreateChatRequest) (*descChat.CreateChatResponse, error) {
	var chatID int64
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	chatQuery := `insert into chats (name) values (@name) returning id;`
	chatQueryArgs := pgx.NamedArgs{
		"name": req.GetName(),
	}

	row := tx.QueryRow(ctx, chatQuery, chatQueryArgs)
	err = row.Scan(&chatID)
	if err != nil {
		return nil, err
	}

	batch := &pgx.Batch{}
	chatUserQuery := `insert into chats_users (chat_id, user_id) values (@chat_id, @user_id);`
	for _, id := range req.GetUserIds() {
		chatUserQueryArgs := pgx.NamedArgs{
			"chat_id": chatID,
			"user_id": id,
		}
		batch.Queue(chatUserQuery, chatUserQueryArgs)
	}

	r := tx.SendBatch(ctx, batch)
	err = r.Close()
	if err != nil {
		return nil, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &descChat.CreateChatResponse{Id: chatID}, nil
}

func (s *server) DeleteChat(ctx context.Context, req *descChat.DeleteChatRequest) (*emptypb.Empty, error) {
	delChat := `delete from chats where id = @id`
	delChatArgs := pgx.NamedArgs{
		"id": req.GetId(),
	}

	_, err := s.pool.Exec(ctx, delChat, delChatArgs)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *descChat.SendMessageRequest) (*emptypb.Empty, error) {
	chatUserQuery := `select id from chats_users where chat_id=@chat_id and user_id=@user_id limit 1;`
	chatUserQueryArgs := pgx.NamedArgs{
		"chat_id": req.GetChatId(),
		"user_id": req.GetUserId(),
	}
	var chatUserID int64
	err := s.pool.QueryRow(ctx, chatUserQuery, chatUserQueryArgs).Scan(&chatUserID)
	if err != nil {
		return nil, err
	}

	msgQuery := `insert into messages (source_id, text) values (@source_id, @text) returning text;`
	msgQueryArgs := pgx.NamedArgs{
		"source_id": chatUserID,
		"text":      req.GetText(),
	}
	_, err = s.pool.Exec(ctx, msgQuery, msgQueryArgs)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
