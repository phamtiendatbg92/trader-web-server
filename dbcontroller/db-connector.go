package dbcontroller

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Tutorial struct {
	Id      primitive.ObjectID `bson:"_id,omitempty"`
	Title   string             `bson:"title,omitempty"`
	Content string             `bson:"content,omitempty"`
	Url     string             `bson:"url,omitempty"`
	Tag     []string           `bson:"tags,omitempty"`
}
type HashTag struct {
	Id   primitive.ObjectID `bson:"_id,omitempty"`
	Tags []string           `bson:"tags,omitempty"`
}

var atlasUri = "mongodb+srv://trader-admin:Q1w2e3r4@cluster0.nwh2i.mongodb.net/trader-web-mongodb?retryWrites=true&w=majority"

var client *mongo.Client

// tutorial mongo collection
var tutCollection *mongo.Collection

// hash tag mongo collection
var hashTagCollection *mongo.Collection

func InitCollection() {
	tutCollection = client.Database("trader-web-mongodb").Collection("tutorials")
	hashTagCollection = client.Database("trader-web-mongodb").Collection("hashtags")
}

func Connect() {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(atlasUri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// Init collection model
	InitCollection()
	//defer client.Disconnect(ctx)
}
