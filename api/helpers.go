package main

import (
	"log"
)

type Auction struct {
	AuctionID    string  `json:"auctionId" bson:"-"`
	StartValue   float64 `json:"startValue" bson:"startValue" binding:"required"`
	BidStartTime int64   `json:"bidStartTime" bson:"bidStartTime" binding:"required"`
	BidEndTime   int64   `json:"bidEndTime" bson:"bidEndTime" binding:"required"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
