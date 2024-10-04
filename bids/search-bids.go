package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func searchBids(params BidSearchParams) []Bid {
	mongoClient := Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	query := bson.M{
		"clientId":  bson.M{"$eq": params.ClientID},
		"auctionId": bson.M{"$eq": params.AuctionID},
	}
	cursor, err := mongoClient.Database("auctions").Collection("bids").Find(ctx, query)
	failOnError(err, "Error fetching bids")
	defer cursor.Close(ctx)

	var bids []Bid
	for cursor.Next(context.Background()) {
		var current Bid
		cursor.Decode(&current)
		bids = append(bids, current)
	}

	return bids
}
