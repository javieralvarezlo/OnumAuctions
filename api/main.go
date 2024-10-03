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
	fmt.Println(response)

	json.Unmarshal(response, &newAuction)

	c.JSON(http.StatusCreated, gin.H{
		"auction": newAuction,
	})
}
