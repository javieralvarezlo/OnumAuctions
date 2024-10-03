package main

import "log"

type Auction struct {
	AuctionID    string  `json:"auctionId"`
	StartValue   float64 `json:"startValue" binding:"required"`
	BidStartTime int64   `json:"bidStartTime" binding:"required"`
	BidEndTime   int64   `json:"bidEndTime" binding:"required"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
