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

func insertBid(bid h.Bid) h.Bid {
	mongoClient := Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	result, err := mongoClient.Database("auctions").Collection("bids").InsertOne(ctx, bid)
	if err != nil {
		fmt.Println("Error inserting bid")
	}

	bid.BidID = result.InsertedID.(primitive.ObjectID).Hex()

	return bid
}

func getBidAuction(bid h.Bid) (h.Auction, error) {
	mongoClient := Connect()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	auctionId, err := primitive.ObjectIDFromHex(bid.AuctionID)
	h.FailOnError(err, "Error converting Object ID from Hex")

	var auction h.Auction

	searchErr := mongoClient.Database("auctions").Collection("auctions").FindOne(ctx, bson.D{{Key: "_id", Value: auctionId}}).Decode(&auction)
	if searchErr != nil {
		return auction, searchErr
	}
	return auction, nil
}

func searchBids(params h.BidSearchParams) []h.Bid {
	mongoClient := Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	query := bson.M{
		"clientId":  bson.M{"$eq": params.ClientID},
		"auctionId": bson.M{"$eq": params.AuctionID},
	}
	cursor, err := mongoClient.Database("auctions").Collection("bids").Find(ctx, query)
	h.FailOnError(err, "Error fetching bids")
	defer cursor.Close(ctx)

	var bids []h.Bid
	for cursor.Next(context.Background()) {
		var current h.Bid
		cursor.Decode(&current)
		bids = append(bids, current)
	}

	return bids
}

func getBestBid(auctionId string) (h.Bid, error) {
	mongoClient := Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	query := bson.D{{Key: "auctionId", Value: auctionId}, {Key: "status", Value: bson.D{{Key: "$eq", Value: "best"}}}}

	var highestBid h.Bid

	err := mongoClient.Database("auctions").Collection("bids").FindOne(ctx, query).Decode(&highestBid)
	if err != nil {
		return highestBid, err
	}

	return highestBid, nil
}
