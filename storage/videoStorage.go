package storage

import "bytes"

type VideoStorage interface {
	Save(videoId string, originalname string, encoding string, mimetype string, videoData bytes.Buffer) (string, error)
}
