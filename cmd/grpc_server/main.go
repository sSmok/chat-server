package main

import (
	"context"
	"log"

	"github.com/sSmok/chat-server/internal/app"
)

func main() {
	ctx := context.Background()
	newApp, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("fail to init app: %v", err)
	}

	err = newApp.Run(ctx)
	if err != nil {
		log.Fatalf("fail to run app: %v", err)
	}
}
