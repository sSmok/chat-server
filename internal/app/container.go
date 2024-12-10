package app

import (
	"context"
	"log"

	chatApi "github.com/sSmok/chat-server/internal/api/chat"
	"github.com/sSmok/chat-server/internal/repository"
	chatRepo "github.com/sSmok/chat-server/internal/repository/chat"
	"github.com/sSmok/chat-server/internal/service"
	chatService "github.com/sSmok/chat-server/internal/service/chat"
	"github.com/sSmok/platform_common/pkg/client/db"
	"github.com/sSmok/platform_common/pkg/client/db/pg"
	"github.com/sSmok/platform_common/pkg/client/db/transaction"
	"github.com/sSmok/platform_common/pkg/config"
)

type container struct {
	pgConfig   config.PGConfigI
	grpcConfig config.GRPCConfigI
	dbClient   db.ClientI
	txManager  db.TxManagerI
	repo       repository.ChatRepositoryI
	serv       service.ChatServiceI
	api        *chatApi.API
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

func (c *container) DBClient(ctx context.Context) db.ClientI {
	if c.dbClient == nil {
		client, err := pg.NewPGClient(ctx, c.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to get pg client: %v", err)
		}
		c.dbClient = client
	}

	return c.dbClient
}

func (c *container) TxManger(ctx context.Context) db.TxManagerI {
	if c.txManager == nil {
		c.txManager = transaction.NewManager(c.DBClient(ctx).DB())
	}

	return c.txManager
}

func (c *container) ChatRepo(ctx context.Context) repository.ChatRepositoryI {
	if c.repo == nil {
		repo := chatRepo.NewChatRepo(c.DBClient(ctx))
		c.repo = repo
	}
	return c.repo
}

func (c *container) ChatService(ctx context.Context) service.ChatServiceI {
	if c.serv == nil {
		serv := chatService.NewChatService(c.ChatRepo(ctx), c.TxManger(ctx))
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
