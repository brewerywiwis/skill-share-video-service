package model

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RawVideoModel struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	VideoID   uuid.UUID          `bson:"video_id" json:"video_id"`
	VideoLink string             `bson:"video_link" json:"video_link"`
	Encoding  string             `bson:"encoding" json:"encoding"`
	Mimetype  string             `bson:"mime_type" json:"mime_type"`
	Size      uint32             `bson:"size" json:"size"`
}
