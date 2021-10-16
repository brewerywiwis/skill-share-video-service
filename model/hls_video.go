package model

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HlsVideoModel struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	VideoID   uuid.UUID          `bson:"video_id" json:"video_id"`
	VideoLink string             `bson:"video_link" json:"video_link"`
}
