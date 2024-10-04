package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func processBid(bid Bid) {
	mongoClient := Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "value", Value: -1}})
	query := bson.D{{Key: "auctionId", Value: bid.AuctionID}, {Key: "status", Value: bson.D{{Key: "$eq", Value: "best"}}}}

	var highestBid Bid

	cursor, err := mongoClient.Database("auctions").Collection("bids").Find(ctx, query, opts)
	failOnError(err, "Error searching bids")
	defer cursor.Close(context.TODO())

	if !cursor.Next(context.TODO()) {
		updateStatus(bid, "best")
		return
	}

	cursor.Decode(&highestBid)

	if bid.Value > highestBid.Value {
		fmt.Println(12213321123)
		fmt.Println(highestBid)
		fmt.Println(highestBid.BidID)
		fmt.Println(highestBid.Status)
		updateStatus(bid, "best")
		updateStatus(highestBid, "outbided")
		return
	}

	updateStatus(bid, "outbided")

}

func updateStatus(bid Bid, status string) {
	mongoClient := Connect()

	collection := mongoClient.Database("auctions").Collection("bids")
	id, _ := primitive.ObjectIDFromHex(bid.BidID)
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := collection.UpdateOne(context.TODO(), filter, update)
	failOnError(err, "Error updating bid")

	bid.Status = status

	defer notifyUser(bid)
}
