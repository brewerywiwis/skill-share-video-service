package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"skillshare/video/config"
	"skillshare/video/database"
	"skillshare/video/mq"
	"skillshare/video/video"

	"google.golang.org/grpc"
)

func main() {
	config.Init()
	database.Init()
	defer database.Disconnect()
	mq.CreateConnection()
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	serverPort := fmt.Sprintf(":%s", port)
	lis, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	v := video.Server{}

	grpcServer := grpc.NewServer()

	video.RegisterChatServiceServer(grpcServer, &v)
	video.RegisterVideoServiceServer(grpcServer, &v)

	log.Printf("Server is running on %s", serverPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
