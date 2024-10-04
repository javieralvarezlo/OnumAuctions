package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createBid(bid Bid) Bid {
	mongoClient := Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	var auction Auction
	auctionId, err := primitive.ObjectIDFromHex(bid.AuctionID)
	failOnError(err, "Error converting Object ID from Hex")

	err = mongoClient.Database("auctions").Collection("auctions").FindOne(ctx, bson.D{{Key: "_id", Value: auctionId}}).Decode(&auction)
	if err != nil {
		bid.BidID = "-1"
		bid.Status = "The selected auction does not exist"
		return bid
	}

	if time.Now().Unix() < auction.BidStartTime {
		bid.BidID = "-1"
		bid.Status = "The auction has not started yet"
		return bid
	}

	if time.Now().Unix() > auction.BidEndTime {
		bid.BidID = "-1"
		bid.Status = "The auction has already finished"
		return bid
	}

	if bid.Value < int(auction.StartValue) {
		bid.BidID = "-1"
		bid.Status = fmt.Sprint("The initial value of the auction is ", auction.StartValue)
		return bid
	}

	bid.Status = "processing"

	result, err := mongoClient.Database("auctions").Collection("bids").InsertOne(ctx, bid)
	if err != nil {
		fmt.Println("Error inserting bid")
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		bid.BidID = oid.Hex()
	}

	defer processBid(bid)

	return bid
}
