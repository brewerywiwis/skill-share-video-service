package video

import (
	"bytes"
	"io"
	"log"
	"skillshare/video/storage"

	"golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Message{Body: "Hello from the server!"}, nil
}

func (server *Server) UploadVideo(stream VideoService_UploadVideoServer) error {
	req, err := stream.Recv()
	if err != nil {
		log.Printf("Upload video rpc error: %s \n", err)
		return err
	}
	originalName := req.GetInfo().GetOriginalname()
	mimetype := req.GetInfo().GetMimetype()
	encoding := req.GetInfo().GetEncoding()
	// size := req.GetInfo().GetSize()
	// log.Printf("RECEIVED:\n name: %s\nmimetype: %s\n encoding: %s\n size: %s\n", originalName, mimetype, encoding, size)

	videoData := bytes.Buffer{}
	videoSize := 0

	for {
		// log.Print("waiting to receive more data")
		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("Received all data")
			break
		}
		if err != nil {
			log.Printf("cannot receive chunk data: %v", err)
			return err
		}

		chunk := req.GetBuffer()
		size := len(chunk)

		videoSize += size

		_, err = videoData.Write(chunk)

		if err != nil {
			log.Printf("cannot write chunk data: %v", err)
			return err
		}
	}

	videoId, result, err := storage.UploadFile(originalName, mimetype, encoding, videoSize, videoData)

	if err != nil {
		log.Printf("cannot upload data: %v", err)
		return err
	}

	res := &VideoUploadResponse{
		Id:   videoId,
		Size: uint32(videoSize),
	}

	err = stream.SendAndClose(res)

	if err != nil {
		log.Printf("cannot send response: %v", err)
		return err
	}

	log.Printf("saved image with id: %s, size: %d", videoId, videoSize)
	log.Printf("Location: %s", result.Location)
	return nil
}
