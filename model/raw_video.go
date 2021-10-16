package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RawVideoModel struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	VideoID   primitive.ObjectID `bson:"video_id" json:"video_id"`
	VideoLink string             `bson:"video_link" json:"video_link"`
	Encoding  string             `bson:"encoding" json:"encoding"`
	Mimetype  string             `bson:"mime_type" json:"mime_type"`
	Size      int                `bson:"size" json:"size"`
	Creator   string             `bson:"creator" json:"creator"`
}
