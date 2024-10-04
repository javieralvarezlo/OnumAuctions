package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func getAllAuctions(params AuctionSearchParams) []Auction {
	mongoClient := Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	query := bson.M{
		"bidStartTime": bson.M{"$gte": params.From},
		"bidEndTime":   bson.M{"$lte": params.To},
	}
	cursor, err := mongoClient.Database("auctions").Collection("auctions").Find(ctx, query)
	failOnError(err, "Error fetching auctions")
	defer cursor.Close(ctx)

	var auctions []Auction
	for cursor.Next(context.Background()) {
		var current Auction
		cursor.Decode(&current)
		auctions = append(auctions, current)
	}

	return auctions
}
