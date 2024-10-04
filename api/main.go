package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {

	router := gin.Default()

	router.POST("/auctions", createAuction)

	router.GET("/auctions", searchAuctions)

	router.POST("/bids", createBid)

	router.Run(":8080")
}

func createAuction(c *gin.Context) {
	var newAuction Auction

	if err := c.ShouldBindJSON(&newAuction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	auctionJson, err := json.Marshal(newAuction)
	if err != nil {
		fmt.Printf("Error marshaling: %s", err)
	}

	response := sendCreationAuction(auctionJson)

	json.Unmarshal(response, &newAuction)

	c.JSON(http.StatusCreated, gin.H{
		"auction": newAuction,
	})
}

func createBid(c *gin.Context) {
	var newBid Bid

	if err := c.ShouldBindJSON(&newBid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bidJson, err := json.Marshal(newBid)
	if err != nil {
		fmt.Printf("Error marshaling: %s", err)
	}

	response := sendCreationBid(bidJson)

	json.Unmarshal(response, &newBid)

	if newBid.BidID == "-1" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": newBid.Status,
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"bid": newBid,
	})
}

func searchAuctions(c *gin.Context) {
	from := c.DefaultQuery("from", fmt.Sprintf("%d", time.Now().Unix()))
	to := c.DefaultQuery("to", fmt.Sprintf("%d", time.Unix(1<<63-62135596801, 999999999).Unix()))

	fromInt, err := strconv.Atoi(from)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "From value is not valid",
		})

		return
	}

	toInt, err := strconv.Atoi(to)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "To value is not valid",
		})

		return
	}

	mongoClient := Connect()

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	query := bson.M{
		"bidStartTime": bson.M{"$gte": fromInt},
		"bidEndTime":   bson.M{"$lte": toInt},
	}
	cursor, err := mongoClient.Database("auctions").Collection("auctions").Find(ctx, query)
	failOnError(err, "Error fetching auctions")
	defer cursor.Close(ctx)

	var auctions []Auction
	for cursor.Next(context.Background()) {
		var current Auction
		fmt.Println(current)
		cursor.Decode(&current)
		fmt.Println(11)
		fmt.Println(current)
		auctions = append(auctions, current)
	}

	if len(auctions) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "There are not auctions in this timeframe",
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"auctions": auctions,
	})

}
