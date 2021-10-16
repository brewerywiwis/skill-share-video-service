package storage

import (
	"bytes"
	"log"
	"skillshare/video/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/google/uuid"
)

func createSession() (*session.Session, error) {
	s3config := config.GetS3Config()
	return session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(s3config.S3_ACCESS_KEY_ID, s3config.S3_SECRET_KEY, ""),
		Region:      aws.String(s3config.S3_REGION)},
	)
}
func UploadFile(originalName string, mimetype string, encoding string, videoSize int, videoData bytes.Buffer) (string, *s3manager.UploadOutput, error) {
	log.Printf("Uploading file")
	session, err := createSession()
	if err != nil {
		log.Println("Cannot create session")
	}
	uploader := s3manager.NewUploader(session)
	config := config.GetS3Config()
	videoId, err := uuid.NewRandom()
	if err != nil {
		log.Printf("cannot random video id in upload video")
		return "", nil, err
	}
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:          aws.String(config.S3_BUCKET),
		Key:             aws.String(config.S3_RAW_VIDEO_KEY + "/" + videoId.String()),
		Body:            bytes.NewReader(videoData.Bytes()),
		ContentType:     aws.String(mimetype),
		ContentEncoding: aws.String(encoding),
	})
	return videoId.String(), result, nil
}
