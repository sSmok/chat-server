package app

import (
	"context"
	"flag"
	"net"

	"github.com/joho/godotenv"
	"github.com/sSmok/chat-server/internal/closer"
	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
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
	app.grpcServer = grpc.NewServer()
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

	if err = app.grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
