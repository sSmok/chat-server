package main

import (
	"context"
	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

const address = "localhost:50501"

func main() {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("connection not created: %v", err)
	}
	defer conn.Close()

	client := descChat.NewChatV1Client(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	//==================
	createReq := &descChat.CreateRequest{
		Info: &descChat.ChatInfo{
			Usernames: []string{"User 1", "User 2", "User 3"},
		},
	}
	createResp, err := client.Create(ctx, createReq)
	if err != nil {
		log.Fatalf("create request failed: %v", err)
	}
	chatId := createResp.GetId()
	log.Printf("chat created successfully: %+v\n", chatId)

	//==================
	_, err = client.Delete(ctx, &descChat.DeleteRequest{Id: chatId})
	if err != nil {
		log.Fatalf("delete request failed: %v", err)
	}
	log.Printf("chat with id=%v deleted successfully", chatId)

	//==================
	msg := &descChat.SendMessageRequest{
		Message: &descChat.Message{
			From:      "User 3",
			Text:      "Chat message text",
			Timestamp: timestamppb.New(time.Now()),
		},
	}
	_, err = client.SendMessage(ctx, msg)
	if err != nil {
		log.Fatalf("message not created: %v", err)
	}
}
