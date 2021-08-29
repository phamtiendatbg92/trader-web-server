package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Tutorial struct {
	//Id string `bson:_id,omitempty`
	Title   string   `bson:"title,omitempty"`
	Content string   `bson:"content,omitempty"`
	Url     string   `bson:"url,omitempty"`
	Tag     []string `bson:"tags,omitempty"`
}

var atlasUri = "mongodb+srv://trader-admin:Q1w2e3r4@cluster0.nwh2i.mongodb.net/trader-web-mongodb?retryWrites=true&w=majority"

var client *mongo.Client

func connect() {
	var err error
	client, err = mongo.NewClient(options.Client().ApplyURI(atlasUri))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	//defer client.Disconnect(ctx)
}
