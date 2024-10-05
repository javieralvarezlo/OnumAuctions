package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	h "auctions/helper"
)

func getAllAuctions(params h.AuctionSearchParams) []h.Auction {
	mongoClient := Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	query := bson.M{
		"bidStartTime": bson.M{"$gte": params.From},
		"bidEndTime":   bson.M{"$lte": params.To},
	}
	cursor, err := mongoClient.Database("auctions").Collection("auctions").Find(ctx, query)
	h.FailOnError(err, "Error fetching auctions")
	defer cursor.Close(ctx)

	var auctions []h.Auction
	for cursor.Next(context.Background()) {
		var current h.Auction
		cursor.Decode(&current)
		auctions = append(auctions, current)
	}

	return auctions
}
