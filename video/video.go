package video

import (
	"log"

	"golang.org/x/net/context"
)

type Server struct{}

func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Message{Body: "Hello from the server!"}, nil
}

func (server *Server) UploadVideo(stream VideoService_UploadVideoServer) error {
	return nil
}
