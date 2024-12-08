package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sSmok/chat-server/internal/config"
	"github.com/sSmok/chat-server/internal/converter"
	repositoryChat "github.com/sSmok/chat-server/internal/repository/chat"
	repositoryMessage "github.com/sSmok/chat-server/internal/repository/message"
	repositoryUser "github.com/sSmok/chat-server/internal/repository/user"
	"github.com/sSmok/chat-server/internal/service"
	"github.com/sSmok/chat-server/internal/service/chat"
	"github.com/sSmok/chat-server/internal/service/message"
	"github.com/sSmok/chat-server/internal/service/user"
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
	userRepo := repositoryUser.NewUserRepo(pool)
	chatRepo := repositoryChat.NewChatRepo(pool)
	messageRepo := repositoryMessage.NewMessageRepo(pool)

	userService := user.NewUserService(userRepo)
	chatService := chat.NewChatService(chatRepo)
	messageService := message.NewMessageService(messageRepo)

	s := &server{
		userService:    userService,
		chatService:    chatService,
		messageService: messageService,
	}

	descChat.RegisterChatV1Server(serv, s)
	if err = serv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %+v\n", err)
	}
}

type server struct {
	userService    service.UserServiceI
	chatService    service.ChatServiceI
	messageService service.MessageServiceI

	descChat.UnimplementedChatV1Server
}

func (s *server) CreateUser(ctx context.Context, req *descChat.CreateUserRequest) (*descChat.CreateUserResponse, error) {
	info := converter.ToUserInfoFromProto(req.GetInfo())
	userID, err := s.userService.CreateUser(ctx, &info)
	if err != nil {
		return nil, err
	}
	return &descChat.CreateUserResponse{Id: userID}, nil
}

func (s *server) DeleteUser(ctx context.Context, req *descChat.DeleteUserRequest) (*emptypb.Empty, error) {
	err := s.userService.DeleteUser(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) CreateChat(ctx context.Context, req *descChat.CreateChatRequest) (*descChat.CreateChatResponse, error) {
	info := converter.ToChatInfoFromProto(req.GetInfo())
	chatID, err := s.chatService.CreateChat(ctx, &info)
	if err != nil {
		return nil, err
	}

	return &descChat.CreateChatResponse{Id: chatID}, nil
}

func (s *server) DeleteChat(ctx context.Context, req *descChat.DeleteChatRequest) (*emptypb.Empty, error) {
	err := s.chatService.DeleteChat(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) CreateMessage(ctx context.Context, req *descChat.CreateMessageRequest) (*emptypb.Empty, error) {
	info := converter.ToMessageInfoFromProto(req.GetInfo())
	err := s.messageService.CreateMessage(ctx, &info)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
