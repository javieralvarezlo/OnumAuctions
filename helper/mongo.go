package helper

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func UpdateStatus(mongoClient *mongo.Client, bid Bid, status string) {
	collection := mongoClient.Database("auctions").Collection("bids")
	id, _ := primitive.ObjectIDFromHex(bid.BidID)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	FailOnError(err, "Error updating bid")

	bid.Status = status

	defer NotifyUser(bid)
}
