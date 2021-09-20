package dbcontroller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FindAllTutorial() ([]Tutorial, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, currErr := tutCollection.Find(ctx, bson.D{})

	var tutorials []Tutorial
	err := cur.All(ctx, &tutorials)
	if err != nil {
		panic(err)
	}
	return tutorials, currErr
}

func CreateNewPost(title string, content string, tags []string, url string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := tutCollection.InsertOne(ctx, bson.D{
		{Key: "title", Value: title},
		{Key: "content", Value: content},
		{Key: "tags", Value: tags},
		{Key: "url", Value: url},
	})
	return err
}
func FindTutByUrl(url string) (Tutorial, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var detailTut Tutorial
	filter := bson.D{bson.E{Key: "url", Value: url}}
	err := tutCollection.FindOne(ctx, filter).Decode(&detailTut)
	return detailTut, err
}

func GetHashTag() (HashTag, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	var hashtag HashTag
	err := hashTagCollection.FindOne(ctx, bson.D{}).Decode(&hashtag)
	if err == mongo.ErrNoDocuments {
		// save tag
	} else if err != nil {
		panic(err)
	}
	return hashtag, err
}
func AddHashTag(tag []string) error {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	_, err := hashTagCollection.InsertOne(ctx, bson.D{{Key: "tags", Value: tag}})
	return err
}
func UpdateHashTag(tag HashTag) error {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	_, err := hashTagCollection.UpdateByID(ctx, tag.Id, bson.D{{Key: "$set", Value: bson.D{{Key: "tags", Value: tag.Tags}}}})
	return err
}

func UpdateTutorial(tut Tutorial) error {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()

	_, err := tutCollection.UpdateByID(ctx, tut.Id, bson.D{{Key: "$set", Value: bson.D{
		{Key: "title", Value: tut.Title},
		{Key: "content", Value: tut.Content},
		{Key: "tags", Value: tut.Tag},
	}}})
	return err
}

func DeleteTutorial(id primitive.ObjectID) error {
	ctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancle()
	_, err := tutCollection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	return err
}
