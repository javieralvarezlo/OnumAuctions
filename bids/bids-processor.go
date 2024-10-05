package main

import (
	h "auctions/helper"
)

func processBid(bid h.Bid) {
	bestBid, err := getBestBid(bid.AuctionID)
	if err != nil {
		//Es la primera bid de la auction
		h.UpdateStatus(Connect(), bid, "best")
		return
	}

	if bid.Value > bestBid.Value {
		h.UpdateStatus(Connect(), bid, "best")
		h.UpdateStatus(Connect(), bestBid, "outbided")
		return
	}

	h.UpdateStatus(Connect(), bid, "outbided")
}