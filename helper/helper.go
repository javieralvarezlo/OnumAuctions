package helper

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func NotifyUser(bid Bid) {
	body, _ := json.Marshal(bid)
	req, _ := http.NewRequest(http.MethodPut, bid.Update, bytes.NewBuffer(body))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	_, err := client.Do(req)
	FailOnError(err, "Error sending the HTTP request")
}

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
