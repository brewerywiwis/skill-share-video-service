package video

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"skillshare/video/config"
	"skillshare/video/model"
	"skillshare/video/mq"
	"skillshare/video/repository"
	"skillshare/video/storage"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
)

type Server struct {
}

type ConverterMessage struct {
	VideoLink string `json:"video_link"`
}

func (s *Server) SayHello(ctx context.Context, in *Message) (*Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &Message{Body: "Hello from the server!"}, nil
}

func (server *Server) UploadVideo(stream VideoService_UploadVideoServer) error {
	res := &VideoUploadResponse{}
	req, err := stream.Recv()
	if err != nil {
		log.Printf("Upload video rpc error: %s \n", err)
		stream.SendAndClose(res)
		return err
	}
	originalName := req.GetInfo().GetOriginalname()
	mimetype := req.GetInfo().GetMimetype()
	encoding := req.GetInfo().GetEncoding()
	creator := req.GetInfo().GetCreator()
	title := req.GetInfo().GetTitle()
	description := req.GetInfo().GetDescription()
	_, err = uuid.Parse(creator)
	if err != nil {
		log.Printf("Creator cannot parse to UUID: %s \n", err)
		stream.SendAndClose(res)
		return err
	}
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
			stream.SendAndClose(res)
			return err
		}

		chunk := req.GetBuffer()
		size := len(chunk)

		videoSize += size

		_, err = videoData.Write(chunk)

		if err != nil {
			log.Printf("cannot write chunk data: %v", err)
			stream.SendAndClose(res)
			return err
		}
	}

	videoId, result, err := storage.UploadFile(originalName, mimetype, encoding, videoSize, videoData)

	if err != nil {
		log.Printf("cannot upload data: %v", err)
		stream.SendAndClose(res)
		return err
	}

	// log.Printf("saved image with id: %s, size: %d", videoId, videoSize)
	// log.Printf("Location: %s", result.Location)
	if err != nil {
		log.Printf("cannot parse uuid: %v", err)
		stream.SendAndClose(res)
		return err
	}
	rawVideo := &model.RawVideoModel{
		ID:          primitive.NewObjectID(),
		VideoID:     videoId,
		VideoLink:   result.Location,
		Encoding:    encoding,
		Mimetype:    mimetype,
		Size:        videoSize,
		Creator:     creator,
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err = repository.CreateRawVideo(rawVideo)
	if err != nil {
		log.Printf("cannot insert result: %v", err)
		storage.DeleteFile(result.Location)
		stream.SendAndClose(res)
		return err
	}

	res = &VideoUploadResponse{
		Id:   videoId.Hex(),
		Size: uint32(videoSize),
	}
	config := config.GetRabbitMQConfig()
	currentQueue := "converter"
	routingKey := currentQueue + config.RoutingKeySuffix
	channel := mq.CreateChannel(currentQueue, routingKey)
	myJsonString, err := json.Marshal(&ConverterMessage{result.Location})
	if err != nil {
		log.Printf("cannot marshal message: %v", err)
		stream.SendAndClose(res)
		return err
	}
	err = mq.Publish(myJsonString, channel)
	if err != nil {
		log.Printf("cannot publish message: %v", err)
		stream.SendAndClose(res)
		return err
	}
	log.Println("[Converter] Send message to converter queue")

	channel.Close()
	log.Println("Send response data")
	err = stream.SendAndClose(res)
	if err != nil {
		log.Printf("cannot send response: %v", err)
		return err
	}

	return nil
}

func (server *Server) GetAllVideo(req *Empty, stream VideoService_GetAllVideoServer) error {
	results, err := repository.GetAllRawVideo()
	if err != nil {
		log.Println("Cannot get all video data from DB")
	}
	for _, v := range *results {
		res := &VideoUploadedResponse{
			VideoUploaded: &VideoUploaded{
				VideoId:     v.VideoID.Hex(),
				VideoLink:   v.VideoLink,
				Title:       v.Title,
				Description: v.Description,
				Creator:     v.Creator,
				CreatedAt:   v.CreatedAt.Format(time.RFC3339),
				UpdatedAt:   v.UpdatedAt.Format(time.RFC3339),
			},
		}
		stream.Send(res)
	}
	return nil
}

func (server *Server) GetRandomVideo(req *NumberRequest, stream VideoService_GetRandomVideoServer) error {
	results, err := repository.GetRandomVideo(int(req.Number))
	if err != nil {
		log.Println("Cannot get random video data from DB")
	}
	for _, v := range *results {
		res := &VideoUploadedResponse{
			VideoUploaded: &VideoUploaded{
				VideoId:     v.VideoID.Hex(),
				VideoLink:   v.VideoLink,
				Title:       v.Title,
				Description: v.Description,
				Creator:     v.Creator,
				CreatedAt:   v.CreatedAt.Format(time.RFC3339),
				UpdatedAt:   v.UpdatedAt.Format(time.RFC3339),
			},
		}
		stream.Send(res)
	}
	return nil
}

func (server *Server) GetVideoByCriteria(req *VideoCriteriaRequest, stream VideoService_GetVideoByCriteriaServer) error {
	filter := map[string]interface{}{}
	if req.GetId() != "" {
		idFromHex, err := primitive.ObjectIDFromHex(req.GetId())
		if err != nil {
			stream.Send(&VideoUploadedResponse{})
			return errors.New("Cannot parse hex to object id")
		}
		filter["video_id"] = idFromHex
	}
	if req.GetUserId() != "" {
		filter["creator"] = req.GetUserId()
	}
	if filter["video_id"] == nil && filter["creator"] == nil {
		stream.Send(&VideoUploadedResponse{})
		return errors.New("Cannot receive filter criteria")
	}
	results, err := repository.GetVideoByCriteria(filter)
	if err != nil {
		log.Println("Cannot get video data from DB")
	}
	for _, v := range *results {
		res := &VideoUploadedResponse{
			VideoUploaded: &VideoUploaded{
				VideoId:     v.VideoID.Hex(),
				VideoLink:   v.VideoLink,
				Title:       v.Title,
				Description: v.Description,
				Creator:     v.Creator,
				CreatedAt:   v.CreatedAt.Format(time.RFC3339),
				UpdatedAt:   v.UpdatedAt.Format(time.RFC3339),
			},
		}
		stream.Send(res)
	}
	return nil
}
