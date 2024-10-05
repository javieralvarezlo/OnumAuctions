package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	h "auctions/helper"
)

func closeAuction(auction h.Auction) {
	mongoClient := Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	query := bson.M{
		"status":    bson.M{"$eq": "outbided"},
		"auctionId": bson.M{"$lte": auction.AuctionID},
	}

	cursor, err := mongoClient.Database("auctions").Collection("bids").Find(ctx, query)
	h.FailOnError(err, "Error fetching auctions")
	defer cursor.Close(ctx)

	for cursor.Next(context.Background()) {
		var current h.Bid
		cursor.Decode(&current)
		h.UpdateStatus(mongoClient, current, "lost")
	}

	query = bson.M{
		"status":    bson.M{"$eq": "best"},
		"auctionId": bson.M{"$lte": auction.AuctionID},
	}

	cursor, err = mongoClient.Database("auctions").Collection("bids").Find(ctx, query)
	h.FailOnError(err, "Error fetching auctions")
	defer cursor.Close(ctx)

	for cursor.Next(context.Background()) {
		var current h.Bid
		cursor.Decode(&current)
		h.UpdateStatus(mongoClient, current, "won")
	}
}
