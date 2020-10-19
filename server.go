package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"papa.com/Clientes/chat"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("Failed to listen on port 9000: %v", err)
	}

	s := chat.Server{}
	grpcServerChat := grpc.NewServer()

	chat.RegisterChatServiceServer(grpcServerChat, &s)

	if err := grpcServerChat.Serve(lis); err != nil {
		log.Fatalf("Failed to server gRPC server over port 9000: %v", err)
	}

}
