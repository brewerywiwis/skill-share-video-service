package storage

import (
	"bytes"
	"log"
	"skillshare/video/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createSession() (*session.Session, error) {
	s3config := config.GetS3Config()
	return session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(s3config.S3_ACCESS_KEY_ID, s3config.S3_SECRET_KEY, ""),
		Region:      aws.String(s3config.S3_REGION)},
	)
}
func UploadFile(originalName string, mimetype string, encoding string, videoSize int, videoData bytes.Buffer) (primitive.ObjectID, *s3manager.UploadOutput, error) {
	log.Printf("Uploading file")
	session, err := createSession()
	if err != nil {
		log.Println("Cannot create session")
	}
	uploader := s3manager.NewUploader(session)
	config := config.GetS3Config()
	videoId := primitive.NewObjectID()
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:          aws.String(config.S3_BUCKET),
		Key:             aws.String(config.S3_RAW_VIDEO_KEY + "/" + videoId.Hex()),
		Body:            bytes.NewReader(videoData.Bytes()),
		ContentType:     aws.String(mimetype),
		ContentEncoding: aws.String(encoding),
	})
	return videoId, result, nil
}

func DeleteFile(path string) error {
	session, err := createSession()
	if err != nil {
		log.Println("Cannot create session")
		return err
	}
	svc := s3.New(session)
	config := config.GetS3Config()
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(config.S3_BUCKET),
		Key:    aws.String(path),
	})
	if err != nil {
		log.Println("Cannot create session")
		return err
	}
	err = svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(config.S3_BUCKET),
		Key:    aws.String(path),
	})
	if err != nil {
		return err
	}
	return nil
}
