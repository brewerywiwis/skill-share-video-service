package main

import (
	"log"
	"net"
	"skillshare/video/video"

	"google.golang.org/grpc"
)

func main() {
	serverPort := ":8100"
	lis, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	v := video.Server{}

	grpcServer := grpc.NewServer()

	video.RegisterChatServiceServer(grpcServer, &v)

	log.Printf("Server is running on %s", serverPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
