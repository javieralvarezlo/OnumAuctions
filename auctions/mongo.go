package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	h "auctions/helper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var once sync.Once

func Connect() *mongo.Client {
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017/"))
		if err != nil {
			log.Fatal("Error connecting to MongoDB:", err)
		}

		mongoClient = client
	})

	return mongoClient
}

func Close() {
	if mongoClient != nil {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Fatal("Error closing MongoDB connection:", err)
		}
	}
}

func insertAuction(auction h.Auction) h.Auction {
	mongoClient := Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := mongoClient.Database("auctions").Collection("auctions").InsertOne(ctx, auction)
	if err != nil {
		fmt.Println("Error inserting auction")
	}

	auction.AuctionID = result.InsertedID.(primitive.ObjectID).Hex()

	return auction
}

func searchAuctions(params h.AuctionSearchParams) []h.Auction {
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

func getAuctionBids(auction h.Auction) []h.Bid {
	var bids []h.Bid
	mongoClient := Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	query := bson.M{
		"auctionId": bson.M{"$eq": auction.AuctionID},
	}

	cursor, err := mongoClient.Database("auctions").Collection("bids").Find(ctx, query)
	h.FailOnError(err, "Error fetching auctions")
	defer cursor.Close(ctx)

	for cursor.Next(context.Background()) {
		var current h.Bid
		cursor.Decode(&current)
		bids = append(bids, current)
	}

	return bids
}
