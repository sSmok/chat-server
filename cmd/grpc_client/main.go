package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/sSmok/chat-server/internal/config"
	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()
	err := config.Load(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := config.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to load GRPC config: %v", err)
	}

	conn, err := grpc.NewClient(grpcConfig.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connection not created: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("connection cannot be closed: %v", err)
		}
	}()

	client := descChat.NewChatV1Client(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// CREATE USER
	createUserReq := &descChat.CreateUserRequest{Info: &descChat.UserInfo{Name: "Nikita"}}
	createUserResp, err := client.CreateUser(ctx, createUserReq)
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	userID := createUserResp.GetId()

	// CREATE CHAT
	createChatReq := &descChat.CreateChatRequest{
		Info: &descChat.ChatInfo{
			Name:    "Chat 1",
			UserIds: []int64{userID},
		},
	}
	createChatResp, err := client.CreateChat(ctx, createChatReq)
	if err != nil {
		log.Fatalf("create request failed: %v", err)
	}
	chatID := createChatResp.GetId()
	log.Printf("chat created successfully: %+v\n", chatID)

	//==================
	msg := &descChat.CreateMessageRequest{
		Info: &descChat.MessageInfo{
			UserId: userID,
			ChatId: chatID,
			Text:   "Text message",
		},
	}
	_, err = client.CreateMessage(ctx, msg)
	if err != nil {
		log.Fatalf("message not created: %v", err)
	}

	// DELETE CHAT
	_, err = client.DeleteChat(ctx, &descChat.DeleteChatRequest{Id: chatID})
	if err != nil {
		log.Fatalf("delete request failed: %v", err)
	}
	log.Printf("chat with id=%v deleted successfully", chatID)
}
