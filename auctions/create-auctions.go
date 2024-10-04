package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createAuction(auction Auction) Auction {
	mongoClient := Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := mongoClient.Database("auctions").Collection("auctions").InsertOne(ctx, auction)
	if err != nil {
		fmt.Println("Error inserting auction")
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		auction.AuctionID = oid.Hex()
	}

	return auction
}
