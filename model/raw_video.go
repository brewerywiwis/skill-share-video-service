package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RawVideoModel struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	VideoID     primitive.ObjectID `bson:"video_id" json:"video_id"`
	VideoLink   string             `bson:"video_link" json:"video_link"`
	Encoding    string             `bson:"encoding" json:"encoding"`
	Mimetype    string             `bson:"mime_type" json:"mime_type"`
	Size        int                `bson:"size" json:"size"`
	Creator     string             `bson:"creator" json:"creator"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
