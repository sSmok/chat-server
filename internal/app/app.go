package app

import (
	"context"
	"flag"
	"log"
	"net"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/joho/godotenv"
	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
	"github.com/sSmok/platform_common/pkg/closer"
	platformInterceptor "github.com/sSmok/platform_common/pkg/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
	flag.Parse()
}

// App представляет основную логику приложения, содержит DI контейнер и сервер gRPC
type App struct {
	container  *container
	grpcServer *grpc.Server
}

// NewApp создает объект приложения
func NewApp(ctx context.Context) (*App, error) {
	app := &App{}
	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

// Run запускает gRPC сервер и контролирует закрытие ресурсов
func (app *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return app.runGRPCServer(ctx)
}

func (app *App) initDeps(ctx context.Context) error {
	deps := []func(context.Context) error{
		app.intiConfig,
		app.initContainer,
		app.initGRPCServer,
	}

	for _, fn := range deps {
		err := fn(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *App) intiConfig(_ context.Context) error {
	err := godotenv.Load(configPath)
	if err != nil {
		return err
	}

	return nil
}

func (app *App) initContainer(_ context.Context) error {
	app.container = newContainer()
	return nil
}

func (app *App) initGRPCServer(ctx context.Context) error {
	chain := grpcMiddleware.ChainUnaryServer(app.container.AccessInterceptor().Access, platformInterceptor.Log)
	app.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(chain))

	reflection.Register(app.grpcServer)
	descChat.RegisterChatV1Server(app.grpcServer, app.container.ChatAPI(ctx))

	return nil
}

func (app *App) runGRPCServer(_ context.Context) error {
	lis, err := net.Listen("tcp", app.container.GRPCConfig().Address())
	if err != nil {
		return err
	}
	closer.Add(lis.Close)

	log.Printf("GRPC server is running on %s", app.container.GRPCConfig().Address())

	if err = app.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
