package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	h "auctions/helper"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.POST("/auctions", createAuction)

	router.GET("/auctions", searchAuctions)

	router.GET("/auctions/:auctionId/bids/:clientId", searchBids)

	router.POST("/bids", createBid)

	router.Run(":8080")
}

func createAuction(c *gin.Context) {
	var newAuction h.Auction

	if err := c.ShouldBindJSON(&newAuction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newAuction.BidEndTime < time.Now().UnixMilli() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The auction can not finish in the past"})
		return
	}

	auctionJson, _ := json.Marshal(newAuction)
	response := h.SendRPC("create_auctions", auctionJson)

	delay := (newAuction.BidEndTime - time.Now().UnixMilli())
	json.Unmarshal(response, &newAuction)

	defer h.SendDelayedMsg(response, delay)

	c.JSON(http.StatusCreated, gin.H{
		"auction": newAuction,
	})
}

func searchAuctions(c *gin.Context) {
	from := c.DefaultQuery("from", fmt.Sprintf("%d", time.Now().UnixMilli()))
	to := c.DefaultQuery("to", fmt.Sprintf("%d", math.MaxInt64))

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

	searchParams := h.AuctionSearchParams{From: fromInt, To: toInt}

	paramsJson, _ := json.Marshal(searchParams)
	response := h.SendRPC("search_auctions", paramsJson)

	var auctions []h.Auction
	json.Unmarshal(response, &auctions)

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

func createBid(c *gin.Context) {
	var newBid h.Bid

	if err := c.ShouldBindJSON(&newBid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bidJson, _ := json.Marshal(newBid)
	response := h.SendRPC("create_bids", bidJson)

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

func searchBids(c *gin.Context) {
	auctionId := c.Param("auctionId")
	clientId := c.Param("clientId")

	searchParams := h.BidSearchParams{ClientID: clientId, AuctionID: auctionId}

	paramsJson, _ := json.Marshal(searchParams)
	response := h.SendRPC("search_bids", paramsJson)

	var bids []h.Bid
	json.Unmarshal(response, &bids)

	if len(bids) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "There are not bids for this client on this auction",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"bids": bids,
	})
}
