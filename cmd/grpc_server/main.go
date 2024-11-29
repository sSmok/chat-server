package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"slices"

	descChat "github.com/sSmok/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

const grpcPort = 50501

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %+v\n", err)
	}
	defer lis.Close()

	serv := grpc.NewServer()
	reflection.Register(serv)
	descChat.RegisterChatV1Server(serv, &server{})
	if err = serv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %+v\n", err)
	}
}

type server struct {
	data []*descChat.Chat
	descChat.UnimplementedChatV1Server
}

func (s *server) Create(_ context.Context, req *descChat.CreateRequest) (*descChat.CreateResponse, error) {
	newChat := &descChat.Chat{
		Id: 1,
		Info: &descChat.ChatInfo{
			Usernames: req.GetInfo().GetUsernames(),
		},
	}
	s.data = append(s.data, newChat)

	return &descChat.CreateResponse{Id: newChat.Id}, nil
}

func (s *server) Delete(_ context.Context, req *descChat.DeleteRequest) (*emptypb.Empty, error) {
	for i, chat := range s.data {
		if chat.Id == req.GetId() {
			s.data = slices.Delete(s.data, i, i+1)
			return &emptypb.Empty{}, nil
		}
	}

	return nil, fmt.Errorf("cant delete chat with id: %v", req.GetId())
}

func (s *server) SendMessage(_ context.Context, req *descChat.SendMessageRequest) (*emptypb.Empty, error) {
	newMsg := req.GetMessage()
	log.Printf("new message: %+v\n", newMsg)

	return &emptypb.Empty{}, nil
}
