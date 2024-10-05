package main

import (
	"fmt"
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
		fmt.Println(22)
		fmt.Printf("EndTime: %d, Now: %d, Diff: %d", newAuction.BidEndTime, time.Now().UnixMilli(), newAuction.BidEndTime-time.Now().UnixMilli())
		c.JSON(http.StatusBadRequest, gin.H{"error": "The auction can not finish in the past"})
		return
	}

	newAuction = sendCreationAuction(newAuction)

	c.JSON(http.StatusCreated, gin.H{
		"auction": newAuction,
	})
}

func createBid(c *gin.Context) {
	var newBid h.Bid

	if err := c.ShouldBindJSON(&newBid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newBid = sendCreationBid(newBid)

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

	searchParams := h.AuctionSearchParams{From: fromInt, To: toInt}

	auctions := sendSearchAuctions(searchParams)

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

func searchBids(c *gin.Context) {
	auctionId := c.Param("auctionId")
	clientId := c.Param("clientId")

	searchParams := h.BidSearchParams{ClientID: clientId, AuctionID: auctionId}

	bids := sendSearchBids(searchParams)

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
