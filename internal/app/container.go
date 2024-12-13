package app

import (
	"context"
	"log"

	"github.com/sSmok/auth/pkg/access_v1"
	chatApi "github.com/sSmok/chat-server/internal/api/chat"
	accessClientPkg "github.com/sSmok/chat-server/internal/client"
	internalConfig "github.com/sSmok/chat-server/internal/config"
	"github.com/sSmok/chat-server/internal/interceptor"
	"github.com/sSmok/chat-server/internal/repository"
	chatRepo "github.com/sSmok/chat-server/internal/repository/chat"
	"github.com/sSmok/chat-server/internal/service"
	chatService "github.com/sSmok/chat-server/internal/service/chat"
	"github.com/sSmok/platform_common/pkg/client/db"
	"github.com/sSmok/platform_common/pkg/client/db/pg"
	"github.com/sSmok/platform_common/pkg/client/db/transaction"
	"github.com/sSmok/platform_common/pkg/closer"
	"github.com/sSmok/platform_common/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type container struct {
	pgConfig          config.PGConfigI
	grpcConfig        config.GRPCConfigI
	accessGRPCConfig  internalConfig.AccessGRPCConfigI
	dbClient          db.ClientI
	txManager         db.TxManagerI
	repo              repository.ChatRepositoryI
	serv              service.ChatServiceI
	accessClient      accessClientPkg.AccessCheckerI
	accessInterceptor interceptor.AccessInterceptorI
	api               *chatApi.API
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

func (c *container) AccessGRPCConfig() internalConfig.AccessGRPCConfigI {
	if c.accessGRPCConfig == nil {
		cfg, err := internalConfig.NewAccessGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get access grpc config: %v", err)
		}
		c.accessGRPCConfig = cfg
	}
	return c.accessGRPCConfig
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

func (c *container) AccessClient() accessClientPkg.AccessCheckerI {
	if c.accessClient == nil {
		conn, err := grpc.NewClient(
			c.AccessGRPCConfig().Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("failed to dial GRPC client: %v", err)
		}
		closer.Add(conn.Close)

		cl := access_v1.NewAccessV1Client(conn)
		accessClient := accessClientPkg.NewAccessClient(cl)
		c.accessClient = accessClient
	}
	return c.accessClient
}

func (c *container) AccessInterceptor() interceptor.AccessInterceptorI {
	if c.accessInterceptor == nil {
		accessInterceptor := interceptor.NewAccessInterceptor(c.AccessClient())
		c.accessInterceptor = accessInterceptor
	}
	return c.accessInterceptor
}

func (c *container) ChatAPI(ctx context.Context) *chatApi.API {
	if c.api == nil {
		api := chatApi.NewAPI(c.ChatService(ctx))
		c.api = api
	}
	return c.api
}
