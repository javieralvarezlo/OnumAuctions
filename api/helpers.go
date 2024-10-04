package main

import (
	"log"
)

type Auction struct {
	AuctionID    string `json:"auctionId" bson:"_id,omitempty"`
	StartValue   int64  `json:"startValue" bson:"startValue" binding:"required"`
	BidStartTime int64  `json:"bidStartTime" bson:"bidStartTime" binding:"required"`
	BidEndTime   int64  `json:"bidEndTime" bson:"bidEndTime" binding:"required"`
}

type Bid struct {
	BidID     string `json:"bidId" bson:"-"`
	AuctionID string `json:"auctionId" bson:"auctionId" binding:"required"`
	Value     int    `json:"value" bson:"value" binding:"required"`
	ClientID  string `json:"clientId" bson:"clientId" binding:"required"`
	Update    string `json:"update" bson:"update" binding:"required"`
	Status    string `json:"status" bson:"status"`
}

type AuctionSearchParams struct {
	From int `json:"from"`
	To   int `json:"to"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
