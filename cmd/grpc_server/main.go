package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sSmok/chat-server/internal/config"
	"github.com/sSmok/chat-server/internal/model"
	"github.com/sSmok/chat-server/internal/repository"
	"github.com/sSmok/chat-server/internal/repository/chat"
	"github.com/sSmok/chat-server/internal/repository/message"
	"github.com/sSmok/chat-server/internal/repository/user"
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
	userRepo := user.NewUserRepo(pool)
	chatRepo := chat.NewChatRepo(pool)
	messageRepo := message.NewMessageRepo(pool)

	s := &server{
		userRepo:    userRepo,
		chatRepo:    chatRepo,
		messageRepo: messageRepo,
	}

	descChat.RegisterChatV1Server(serv, s)
	if err = serv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %+v\n", err)
	}
}

type server struct {
	userRepo    repository.UserRepositoryI
	chatRepo    repository.ChatRepositoryI
	messageRepo repository.MessageRepositoryI
	descChat.UnimplementedChatV1Server
}

func (s *server) CreateUser(ctx context.Context, req *descChat.CreateUserRequest) (*descChat.CreateUserResponse, error) {
	userID, err := s.userRepo.CreateUser(ctx, &model.UserInfo{Name: req.GetName()})
	if err != nil {
		return nil, err
	}
	return &descChat.CreateUserResponse{Id: userID}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *descChat.DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.userRepo.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) CreateChat(ctx context.Context, req *descChat.CreateChatRequest) (*descChat.CreateChatResponse, error) {
	chatID, err := s.chatRepo.CreateChat(ctx, req.GetName(), req.GetUserIds())
	if err != nil {
		return nil, err
	}

	return &descChat.CreateChatResponse{Id: chatID}, nil
}

func (s *server) DeleteChat(ctx context.Context, req *descChat.DeleteChatRequest) (*emptypb.Empty, error) {
	err := s.chatRepo.DeleteChat(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *descChat.SendMessageRequest) (*emptypb.Empty, error) {
	err := s.messageRepo.CreateMessage(ctx, req.GetChatId(), req.GetUserId(), req.GetText())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
