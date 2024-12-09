package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	chatApi "github.com/sSmok/chat-server/internal/api/chat"
	"github.com/sSmok/chat-server/internal/closer"
	"github.com/sSmok/chat-server/internal/config"
	"github.com/sSmok/chat-server/internal/repository"
	chatRepo "github.com/sSmok/chat-server/internal/repository/chat"
	"github.com/sSmok/chat-server/internal/service"
	chatService "github.com/sSmok/chat-server/internal/service/chat"
)

type container struct {
	pgConfig   config.PGConfigI
	grpcConfig config.GRPCConfigI

	pool *pgxpool.Pool

	repo repository.ChatRepositoryI
	serv service.ChatServiceI
	api  *chatApi.API
}

func newContainer() *container {
	return &container{}
}

func (c *container) PGConfig() config.PGConfigI {
	if c.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %v", err)
		}
		c.pgConfig = cfg
	}
	return c.pgConfig
}

func (c *container) GRPCConfig() config.GRPCConfigI {
	if c.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %v", err)
		}
		c.grpcConfig = cfg
	}
	return c.grpcConfig
}

func (c *container) Pool(ctx context.Context) *pgxpool.Pool {
	if c.pool == nil {
		pool, err := pgxpool.New(ctx, c.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to get pg pool: %v", err)
		}
		closer.Add(func() error {
			pool.Close()
			return nil
		})

		c.pool = pool
	}
	return c.pool
}

func (c *container) ChatRepo(ctx context.Context) repository.ChatRepositoryI {
	if c.repo == nil {
		repo := chatRepo.NewChatRepo(c.Pool(ctx))
		c.repo = repo
	}
	return c.repo
}

func (c *container) ChatService(ctx context.Context) service.ChatServiceI {
	if c.serv == nil {
		serv := chatService.NewChatService(c.ChatRepo(ctx))
		c.serv = serv
	}
	return c.serv
}

func (c *container) ChatAPI(ctx context.Context) *chatApi.API {
	if c.api == nil {
		api := chatApi.NewAPI(c.ChatService(ctx))
		c.api = api
	}
	return c.api
}
