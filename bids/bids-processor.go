package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	h "auctions/helper"
)

func processBid(bid h.Bid) {
	mongoClient := Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "value", Value: -1}})
	query := bson.D{{Key: "auctionId", Value: bid.AuctionID}, {Key: "status", Value: bson.D{{Key: "$eq", Value: "best"}}}}

	var highestBid h.Bid

	cursor, err := mongoClient.Database("auctions").Collection("bids").Find(ctx, query, opts)
	h.FailOnError(err, "Error searching bids")
	defer cursor.Close(context.TODO())

	if !cursor.Next(context.TODO()) {
		h.UpdateStatus(mongoClient, bid, "best")
		return
	}

	cursor.Decode(&highestBid)

	if bid.Value > highestBid.Value {
		h.UpdateStatus(mongoClient, bid, "best")
		h.UpdateStatus(mongoClient, highestBid, "outbided")
		return
	}

	h.UpdateStatus(mongoClient, bid, "outbided")

}
