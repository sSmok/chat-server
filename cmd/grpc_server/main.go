package main

import (
	"context"
	"flag"
	"log"
	"net"

	"github.com/jackc/pgx/v5/pgxpool"
	chatApi "github.com/sSmok/chat-server/internal/api/chat"
	"github.com/sSmok/chat-server/internal/config"
	repositoryChat "github.com/sSmok/chat-server/internal/repository/chat"
	"github.com/sSmok/chat-server/internal/service/chat"
	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	chatRepo := repositoryChat.NewChatRepo(pool)
	chatService := chat.NewChatService(chatRepo)

	s := chatApi.NewAPI(chatService)

	descChat.RegisterChatV1Server(serv, s)
	if err = serv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %+v\n", err)
	}
}
