package repository

import (
	"context"
	"log"
	"skillshare/video/config"
	"skillshare/video/database"
	"skillshare/video/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func logError(err error) {
	log.Printf("Mongo error %s", err)
}
func CreateRawVideo(data *model.RawVideoModel) (*mongo.InsertOneResult, error) {
	databaseConfig := config.GetDatabaseConfig()
	client := database.GetDatabaseClient()
	collection := client.Database(databaseConfig.DB_NAME).Collection("raw_videos")
	var session mongo.Session
	var err error
	var ctx = context.Background()
	var result *mongo.InsertOneResult
	if session, err = client.StartSession(); err != nil {
		logError(err)
		return nil, err
	}
	if err = session.StartTransaction(); err != nil {
		logError(err)
		return nil, err
	}
	if err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		_, err = collection.InsertOne(sc, data)
		if err != nil {
			logError(err)
			return err
		}
		if err = session.CommitTransaction(sc); err != nil {
			logError(err)
			return err
		}
		return nil
	}); err != nil {
		logError(err)
		return nil, err
	}

	session.EndSession(ctx)
	return result, nil
}

func GetAllRawVideo() (*[]model.RawVideoModel, error) {
	databaseConfig := config.GetDatabaseConfig()
	collection := database.GetDatabaseClient().Database(databaseConfig.DB_NAME).Collection("raw_videos")

	filter := bson.M{}
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Println("Mongo error", err)
		return nil, err
	}

	var results []model.RawVideoModel
	err = cur.All(context.TODO(), &results)
	if err != nil {
		log.Println("Mongo error", err)
		return nil, err
	}
	return &results, nil

}

func UpdateRawVideo() {

}
func DeleteRawVideo() {

}

func GetRandomVideo(n int) (*[]model.RawVideoModel, error) {
	databaseConfig := config.GetDatabaseConfig()
	collection := database.GetDatabaseClient().Database(databaseConfig.DB_NAME).Collection("raw_videos")

	pipeline := mongo.Pipeline{
		bson.D{{"$sample", bson.D{{"size", n}}}},
	}
	cur, err := collection.Aggregate(context.TODO(), pipeline)

	var results []model.RawVideoModel
	err = cur.All(context.TODO(), &results)
	if err != nil {
		log.Println("Mongo error", err)
		return nil, err
	}
	return &results, nil
}

func GetVideoById(videoId string) (*[]model.RawVideoModel, error) {
	databaseConfig := config.GetDatabaseConfig()
	collection := database.GetDatabaseClient().Database(databaseConfig.DB_NAME).Collection("raw_videos")

	videoIdHex, err := primitive.ObjectIDFromHex(videoId)
	if err != nil {
		log.Println("video id error", err)
	}
	filter := bson.D{{"video_id", videoIdHex}}
	cur, err := collection.Find(context.TODO(), filter)

	var results []model.RawVideoModel
	err = cur.All(context.TODO(), &results)
	if err != nil {
		log.Println("Mongo error", err)
		return nil, err
	}
	return &results, nil
}
