package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.POST("/auctions", createAuction)

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
